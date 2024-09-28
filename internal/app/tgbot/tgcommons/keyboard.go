package tgcommons

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (bot *Bot) SendStartUpKeyboard(chatID int64) error {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("База знаний"),
			tgbotapi.NewKeyboardButton("TRS Solver"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Выберите опцию:")
	msg.ReplyMarkup = keyboard

	_, err := bot.bot.Send(msg)

	return err
}

/*
func (bot *TGBotApp) deleteCommand(update tgbotapi.Update) error {
	err := bot.removeKeyboard(update.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("failed to remove keyboard: %w", err)
	}

	return nil
}
*/

func (bot *Bot) RemoveKeyboard(chatID int64, messageText string) error {
	removeKeyboard := tgbotapi.NewRemoveKeyboard(true)

	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = removeKeyboard

	_, err := bot.bot.Send(msg)

	return err
}
