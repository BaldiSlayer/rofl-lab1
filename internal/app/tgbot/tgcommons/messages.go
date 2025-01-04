package tgcommons

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) SendMessageWithReturningID(chatID int64, messageText string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, messageText)

	message, err := bot.bot.Send(msg)

	return message.MessageID, err
}

func (bot *Bot) SendMessage(chatID int64, messageText string) error {
	msg := tgbotapi.NewMessage(chatID, messageText)

	_, err := bot.bot.Send(msg)

	return err
}

func (bot *Bot) SendMarkdownMessage(chatID int64, messageText string) error {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.bot.Send(msg)

	return err
}

func (bot *Bot) EditMarkdownMessage(chatID int64, messageID int, messageText string) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, messageText)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) SendMarkdownMessageWithKeyboard(chatID int64, messageText string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.bot.Send(msg)

	return err
}
