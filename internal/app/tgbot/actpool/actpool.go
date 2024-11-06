package actpool

import (
	"errors"
	"fmt"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/ustorage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateTransition = func(update tgbotapi.Update) (models.UserState, error)

type BotActionsPool struct {
	storage ustorage.UserDataStorage

	actions  map[models.UserState]StateTransition
	commands map[string]StateTransition
}

func New(storage ustorage.UserDataStorage,
	transitions map[models.UserState]StateTransition,
	commands map[string]StateTransition) (*BotActionsPool, error) {
	return &BotActionsPool{
		storage:  storage,
		actions:  transitions,
		commands: commands,
	}, nil
}

// Exec находит для юзера его текущий стейт и исполняет соответствующую функцию
func (b *BotActionsPool) Exec(update tgbotapi.Update) error {
	if update.Message != nil && update.Message.IsCommand() {
		return b.ExecCommand(update)
	}

	userID := update.SentFrom().ID

	userState, err := getUserState(userID, b.storage)
	if err != nil {
		return err
	}

	f, ok := b.actions[userState]
	if !ok {
		return fmt.Errorf("action pooler has no action for this state: %v", userState)
	}

	currentState, err := f(update)
	if err != nil {
		currentState = models.GetRequest
	}

	return errors.Join(err, b.storage.SetState(userID, currentState))
}

func (b *BotActionsPool) ExecCommand(update tgbotapi.Update) error {
	f, ok := b.commands[update.Message.Command()]
	if !ok {
		return fmt.Errorf("action pooler has no action for this command: %s", update.Message.Command())
	}

	state, err := f(update)
	if err != nil {
		state = models.GetRequest
	}

	return errors.Join(err, b.storage.SetState(update.SentFrom().ID, state))
}

func getUserState(userID int64, storage ustorage.UserDataStorage) (models.UserState, error) {
	userState, err := storage.GetState(userID)
	if errors.Is(err, ustorage.ErrNotFound) {
		userState = models.Start
		err = storage.SetState(userID, userState)
	}
	return userState, err
}
