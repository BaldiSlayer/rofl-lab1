package tgbot

import (
	"log/slog"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/mclient"
	appmodels "github.com/BaldiSlayer/rofl-lab1/internal/app/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/actpool"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/controllers"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgcommons"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgconfig"
)

type App struct {
	config *tgconfig.TGBotConfig
	bot    *tgcommons.Bot

	actionsPooler *actpool.BotActionsPool
	userLocks map[int64]*sync.Mutex
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

func New(opts ...Option) *App {
	tgBot := &App{}

	for _, opt := range opts {
		err := opt(tgBot)
		if err != nil {
			slog.Error("failed to init telegram bot", "error", err)
			os.Exit(1)
		}
	}

	var err error

	tgBot.bot, err = tgcommons.NewBot(tgBot.config.Token)
	if err != nil {
		slog.Error("error while creating new bot api instance", "error", err)
		// FIXME: возвращать ошибку, а не выходить
		os.Exit(1)
	}

	tgBot.actionsPooler, err = actpool.New(/*TODO: pass transitions here*/)
	if err != nil {
		slog.Error("error while creating new actpool", "error", err)
		os.Exit(1)
	}

	// TODO: remove
	err = tgBot.initControllers()
	if err != nil {
		slog.Error("error initializing controllers", "error", err)
		os.Exit(1)
	}

	slog.Debug("initialized tgbot api")

	return tgBot
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
		go func(update tgbotapi.Update) {
			userID := update.SentFrom().ID
			userLock := bot.lockByUserID(userID)
			defer userLock.Unlock()

			err := bot.actionsPooler.Exec(update)
			if err != nil {
				slog.Error("failed to process user action", "error", err)
				return
			}

			if update.CallbackQuery != nil {
				err := bot.bot.SendCallbackResponse(update)
				if err != nil {
					slog.Error("failed to send callback response", "error", err)
				}
			}
		}(update)
	}
}

func (bot *App) initControllers() error {
	context, err := appmodels.LoadQABase()
	if err != nil {
		return err
	}

	mclient, err := mclient.NewMistralClient(context)
	if err != nil {
		return err
	}

	controller := controllers.Controller{
		Bot:         bot.bot,
		ModelClient: mclient,
	}

	bot.actionsPooler.AddStateTransition(
		models.Start,
		controller.Start,
	)

	bot.actionsPooler.AddStateTransition(
		models.GetRequest,
		controller.GetRequest,
	)

	bot.actionsPooler.AddStateTransition(
		models.GetTrs,
		controller.GetTrs,
	)

	bot.actionsPooler.AddStateTransition(
		models.ValidateTrs,
		controller.ValidateTrs,
	)

	bot.actionsPooler.AddStateTransition(
		models.FixTrs,
		controller.FixTrs,
	)

	return nil
}

func (bot *App) lockByUserID(userID int64) *sync.Mutex {
	if lock, ok := bot.userLocks[userID]; ok {
		lock.Lock()
		return lock
	}

	lock := &sync.Mutex{}
	lock.Lock()
	bot.userLocks[userID] = lock
	return lock
}
