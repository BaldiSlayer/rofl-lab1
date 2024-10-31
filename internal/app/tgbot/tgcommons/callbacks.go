package tgcommons

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (bot *Bot) SendCallbackResponse(update tgbotapi.Update) error {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	_, err := bot.bot.Request(callback)
	return err
}
