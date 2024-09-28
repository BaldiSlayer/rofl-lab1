package tgcommons

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) SendMessage(chatID int64, messageText string) error {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	msg := tgbotapi.NewMessage(chatID, messageText)

	_, err := bot.bot.Send(msg)

	return err
}
