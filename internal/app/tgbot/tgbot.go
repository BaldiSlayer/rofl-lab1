package tgbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

func New(ctx context.Context, opts ...Option) (*App, error) {
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

	tgBot.bot, err = tgcommons.NewBot(tgBot.config.TgToken)
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

	func() {
		var wg sync.WaitGroup
		for {
			select {
			case update := <-updates:
				slog.Debug("processing update")
				wg.Add(1)
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
					defer cancel()

					bot.processUpdate(ctx, update)
					wg.Done()
				}()
			case <-ctx.Done():
				slog.Info("Gracefully shutting down")
				bot.bot.StopReceivingUpdates()

				wg.Wait()
				bot.ustorageCloser.Close()
				return
			}
		}
	}()
}

func buildTransitions(controller *controllers.Controller) map[models.UserState]actpool.StateTransition {
	return map[models.UserState]actpool.StateTransition{
		models.Start:       controller.Start,
		models.GetRequest:  controller.GetRequest,
		models.GetTrs:      controller.GetTrs,
		models.ValidateTrs: controller.ValidateTrs,
		models.FixTrs:      controller.FixTrs,
	}
}

func buildCommands(controller *controllers.Controller) map[string]actpool.StateTransition {
	return map[string]actpool.StateTransition{
		"start":   controller.StartCommand,
		"help":    controller.HelpCommand,
		"trs":     controller.TrsCommand,
		"version": controller.VersionCommand,
	}
}

func (bot *App) processUpdate(ctx context.Context, update tgbotapi.Update) {
	userID := update.SentFrom().ID
	// TODO: userLock := bot.lockByUserID(userID)
	// if !userLock.TryLock() {
	// 	err := bot.bot.SendMessage(userID, "Предыдущее сообщение еще обрабатывается")
	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 	}
	// 	return
	// }
	// defer userLock.Unlock()

	err := bot.actionsPooler.Exec(ctx, update)
	if err != nil {
		err = errors.Join(err, bot.bot.SendMessage(userID, fmt.Sprintf("Ошибка при обработке запроса: %s", err)))
		err = errors.Join(err, bot.bot.SendMessage(userID, "Введите запрос к Базе Знаний"))
		slog.Error("failed to process user action", "error", err)
		return
	}

	if update.CallbackQuery != nil {
		err := bot.bot.SendCallbackResponse(update)
		if err != nil {
			slog.Error("failed to send callback response", "error", err)
		}
	}
}
