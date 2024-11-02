package actpool

import (
	"fmt"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/ustorage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateTransition = func(update tgbotapi.Update) (models.UserState, error)

type BotActionsPool struct {
	storage ustorage.UserDataStorage

	actions map[models.UserState]StateTransition
}

func New(storage ustorage.UserDataStorage, transitions map[models.UserState]StateTransition) (*BotActionsPool, error) {
	return &BotActionsPool{
		actions: transitions,
		storage: storage,
	}, nil
}

// Exec находит для юзера его текущий стейт и исполняет соответствующую функцию
func (b *BotActionsPool) Exec(update tgbotapi.Update) error {
	userState, err := b.storage.GetState(update.SentFrom().ID)
	if err != nil {
		return err
	}

	f, ok := b.actions[userState]
	if !ok {
		return fmt.Errorf("action pooler has no action for this state: %v", userState)
	}

	currentState, err := f(update)
	if err != nil {
		return err
	}

	err = b.storage.SetState(update.SentFrom().ID, currentState)
	if err != nil {
		return err
	}

	return err
}
