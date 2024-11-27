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

	wh.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}
	wh.MaxConnections = 40

	_, err = bot.Request(wh)
	if err != nil {
		return nil, err
	}

	_, err = bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	return &Bot{
		bot: bot,
	}, err
}

func (bot *Bot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	updates := bot.bot.ListenForWebhook("/" + bot.bot.Token)
	go http.ListenAndServe("0.0.0.0:8443", nil)

	return updates
}

func (bot *Bot) StopReceivingUpdates() {
	bot.bot.StopReceivingUpdates()
}
