package tgbot

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"log/slog"
	"os"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/actpool"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/beclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/controllers"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgcommons"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/tgconfig"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	config *tgconfig.TGBotConfig
	bot    *tgcommons.Bot

	actionsPooler *actpool.BotActionsPool
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
		os.Exit(1)
	}

	tgBot.actionsPooler, err = actpool.New()
	if err != nil {
		slog.Error("error while creating new actpool", "error", err)
		os.Exit(1)
	}

	err = tgBot.initControllers()

	return tgBot
}

func (bot *App) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.bot.GetUpdatesChan(u)

	slog.Info("telegram bot has successfully started")

	for update := range updates {
		// походу надо будет сюда еще воткнуть rate limiter
		go func(update tgbotapi.Update) {
			err := bot.actionsPooler.Exec(update)
			if err != nil {
				slog.Error("failed to process user action", "error", err)
			}
		}(update)
	}
}

func (bot *App) initControllers() error {
	controller := controllers.Controller{
		Bot:           bot.bot,
		BackendClient: &beclient.MockBackendClient{},
	}

	bot.actionsPooler.AddController(
		models.EmptyState,
		controller.EmptyState,
	)

	bot.actionsPooler.AddController(
		models.StartState,
		controller.Start,
	)

	bot.actionsPooler.AddController(
		models.WaitForKBQuestion,
		controller.WaitForKBQuestion,
	)

	return nil
}
