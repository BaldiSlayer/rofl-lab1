package tgcommons

import (
	"log/slog"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

const (
	webhookBaseUrl = "https://tfl-lab1.starovoytovai.ru/"
)

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	wh, err := tgbotapi.NewWebhook(webhookBaseUrl + token)
	if err != nil {
		return nil, err
	}

	_, err = bot.Request(wh)
	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		slog.Error("failed to get webhook info", "error", err)
	}

	if info.LastErrorDate != 0 {
		slog.Error("telegram callback failed", "error", info.LastErrorMessage)
	}

	return &Bot{
		bot: bot,
	}, err
}

func (bot *Bot) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	updates := bot.bot.ListenForWebhook("/" + bot.bot.Token)
	go http.ListenAndServe("0.0.0.0:8443", nil)

	return updates
}

func (bot *Bot) StopReceivingUpdates() {
	bot.bot.StopReceivingUpdates()
}
