package tgcommons

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) SendMessage(chatID int64, messageText string) error {
	msg := tgbotapi.NewMessage(chatID, messageText)

	_, err := bot.bot.Send(msg)

	return err
}

func (bot *Bot) SendMarkdownMessage(chatID int64, messageText string) error {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := bot.bot.Send(msg)

	return err
}

func (bot *Bot) SendMarkdownMessageWithKeyboard(chatID int64, messageText string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := bot.bot.Send(msg)

	return err
}
