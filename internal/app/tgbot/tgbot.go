package tgbot

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

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

	actionsPooler *actpool.BotActionsPool
	userLocks     map[int64]*sync.Mutex
	userStorage   ustorage.UserDataStorage
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

func New(opts ...Option) (*App, error) {
	tgBot := &App{}

	for _, opt := range opts {
		err := opt(tgBot)
		if err != nil {
			return nil, err
		}
	}

	tgBot.userLocks = make(map[int64]*sync.Mutex)

	userStorage, err := ustorage.NewPostgresUserStorage()
	if err != nil {
		return nil, err
	}

	tgBot.userStorage = userStorage

	tgBot.bot, err = tgcommons.NewBot(tgBot.config.Token)
	if err != nil {
		return nil, err
	}

	controller, err := controllers.New(tgBot.bot, userStorage)
	if err != nil {
		return nil, err
	}

	tgBot.actionsPooler, err = actpool.New(userStorage, buildTransitions(controller))
	if err != nil {
		return nil, err
	}

	slog.Debug("initialized tgbot api")

	return tgBot, nil
}

func (bot *App) Run() {
	// NOTE: Offset value set to 0 means that when backend is restarted, updates
	// received by the last call to getUpdates will be resent by the Telegram
	// API, whether they're already handled or not.
	u := tgbotapi.NewUpdate(0)
	// NOTE: Updates per request
	u.Limit = 100
	// NOTE: Timeout of long polling requests
	u.Timeout = 1
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := bot.bot.GetUpdatesChan(u)

	slog.Info("telegram bot has successfully started")

	for update := range updates {
		slog.Debug("processing update")
		go bot.processUpdate(update)
	}
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

func (bot *App) lockByUserID(userID int64) *sync.Mutex {
	if lock, ok := bot.userLocks[userID]; ok {
		return lock
	}

	lock := &sync.Mutex{}
	bot.userLocks[userID] = lock
	return lock
}

func (bot *App) processUpdate(update tgbotapi.Update) {
	userID := update.SentFrom().ID
	userLock := bot.lockByUserID(userID)
	if !userLock.TryLock() {
		err := bot.bot.SendMessage(userID, "Предыдущее сообщение еще обрабатывается")
		if err != nil {
			slog.Error(err.Error())
		}
		return
	}
	defer userLock.Unlock()

	err := bot.actionsPooler.Exec(update)
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
