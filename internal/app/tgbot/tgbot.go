package tgbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/actpool"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/controllers"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgcommons"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgconfig"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/ustorage"
)

type App struct {
	config *tgconfig.TGBotConfig
	bot    *tgcommons.Bot

	controller *controllers.Controller

	actionsPooler  *actpool.BotActionsPool
	userLocks      ustorage.UserLockStorage
	userStorage    ustorage.UserDataStorage
	ustorageCloser ustorage.Closer
}

type Option func(options *App) error

// WithConfig инициализирует конфиг
func WithConfig() Option {
	return func(options *App) error {
		config, err := tgconfig.LoadTGBotConfig()
		if err != nil {
			return err
		}

		options.config = config

		return nil
	}
}

func New(ctx context.Context, callbackMode bool, opts ...Option) (*App, error) {
	tgBot := &App{}

	for _, opt := range opts {
		err := opt(tgBot)
		if err != nil {
			return nil, err
		}
	}

	postgres, err := ustorage.NewPostgresStorage(tgBot.config.PgUser, tgBot.config.PgPassword, tgBot.config.PgDBName)
	if err != nil {
		return nil, err
	}

	tgBot.userLocks = postgres
	tgBot.ustorageCloser = postgres
	tgBot.userStorage = postgres

	tgBot.bot, err = tgcommons.NewBot(tgBot.config.TgToken, callbackMode)
	if err != nil {
		return nil, err
	}

	controller, err := controllers.New(tgBot.bot, postgres, tgBot.config.GhToken)
	if err != nil {
		return nil, err
	}
	tgBot.controller = controller

	tgBot.actionsPooler, err = actpool.New(postgres, buildTransitions(controller), buildCommands(controller))
	if err != nil {
		return nil, err
	}

	slog.Debug("initialized tgbot api")

	return tgBot, nil
}

func (bot *App) Run(ctx context.Context) {
	updates := bot.bot.GetUpdatesChan()

	slog.Info("telegram bot has successfully started")

	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			slog.Info("Gracefully shutting down")

			bot.bot.StopReceivingUpdates()
			wg.Wait()

			bot.ustorageCloser.Close()
			return
		case update := <-updates:
			slog.Debug("processing update")
			wg.Add(1)
			go func() {
				defer wg.Done()

				timeout := time.Minute * 6
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				err := bot.processUpdate(ctx, cancel, timeout, update)
				if err != nil {
					slog.Error("error while processing update", "error", err)
				}
			}()
		}
	}
}

func (bot *App) processUpdate(ctx context.Context, cancel context.CancelFunc, timeout time.Duration, update tgbotapi.Update) error {
	defer bot.processCallbackQuery(update)

	if update.Message != nil && update.Message.Command() == "cancel" {
		bot.processCancelCommand(ctx, update)
		return nil
	}

	userID := update.SentFrom().ID
	requestID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	ok, err := bot.userLocks.TryLock(ctx, userID, requestID, timeout)
	if err != nil {
		return fmt.Errorf("failed to acquire user lock: %w", err)
	}
	if !ok {
		err := bot.bot.SendMessage(userID, "Предыдущее сообщение еще обрабатывается\n\nВы можете отменить его обработку командой /cancel")
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return nil
	}

	defer func() {
		err := bot.userLocks.Unlock(context.Background(), userID, requestID)
		if err != nil {
			slog.Error("failed to unlock user", "error", err)
		}
	}()

	go bot.cancelIfNotLocked(ctx, cancel, userID, requestID)

	state, err := bot.actionsPooler.Exec(ctx, update)
	if errors.Is(ctx.Err(), context.Canceled) {
		err = errors.Join(
			ctx.Err(),
			bot.bot.SendMessage(userID, "Запрос отменен"),
		)
		slog.Info("request cancelled", "error", err)
		return nil
	}
	if err != nil {
		state = models.GetRequest
	}

	err = errors.Join(err, bot.userStorage.SetState(ctx, userID, state))
	if err != nil {
		err = errors.Join(
			err,
			bot.userStorage.SetState(ctx, userID, models.GetRequest),
			bot.bot.SendMessage(userID, fmt.Sprintf("Ошибка при обработке запроса: %s", err)),
			bot.bot.SendMessage(userID, "Введите запрос к Базе Знаний"),
		)
		return fmt.Errorf("failed to process user action: %w", err)
	}
	return nil
}

func (bot *App) processCancelCommand(ctx context.Context, update tgbotapi.Update) {
	userID := update.SentFrom().ID

	err := bot.userLocks.ForceUnlock(ctx, userID)
	if err != nil {
		slog.Error("failed to force unlock user", "userID", userID, "error", err)
	}
}

func (bot *App) processCallbackQuery(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	err := bot.bot.SendCallbackResponse(update)
	if err != nil {
		slog.Error("failed to send callback response", "error", err)
	}
}

func (bot *App) cancelIfNotLocked(ctx context.Context, cancel context.CancelFunc, userID int64, requestID uuid.UUID) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			isLocked, err := bot.userLocks.IsLocked(ctx, userID, requestID)
			if err != nil {
				slog.Error("failed to check user lock", "error", err)
			}
			if !isLocked || err != nil {
				cancel()
				return
			}
		}
	}
}

func buildTransitions(controller *controllers.Controller) map[models.UserState]actpool.StateTransition {
	return map[models.UserState]actpool.StateTransition{
		models.Start:                  controller.Start,
		models.GetRequest:             controller.GetRequest,
		models.GetTrs:                 controller.GetTrs,
		models.ValidateTrs:            controller.ValidateTrs,
		models.FixTrs:                 controller.FixTrs,
		models.GetQuestionMultiModels: controller.GetRequestMultiModels,
		models.GetSimilar:             controller.GetSimilar,
	}
}

func buildCommands(controller *controllers.Controller) map[string]actpool.StateTransition {
	return map[string]actpool.StateTransition{
		"start":       controller.StartCommand,
		"help":        controller.HelpCommand,
		"trs":         controller.TrsCommand,
		"version":     controller.VersionCommand,
		"multimodels": controller.CommandMultiModels,
		"similar":     controller.CommandGetSimilar,
	}
}
