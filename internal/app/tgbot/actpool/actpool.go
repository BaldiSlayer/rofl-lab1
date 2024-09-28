package actpool

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/models"
	"sync"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/tgbot/ustorage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrNoAction = fmt.Errorf("action pooler has no action for this state")

type BotActionsPool struct {
	storage ustorage.UserDataStorage

	actions      map[models.UserState]func(update tgbotapi.Update) (models.UserState, error)
	actionsMutex sync.Mutex
}

func New() (*BotActionsPool, error) {
	storage, err := ustorage.NewMapUserStorage()
	if err != nil {
		return nil, err
	}

	return &BotActionsPool{
		actions: make(map[models.UserState]func(update tgbotapi.Update) (models.UserState, error)),
		storage: storage,
	}, nil
}

func (b *BotActionsPool) AddController(state models.UserState, f func(update tgbotapi.Update) (models.UserState, error)) {
	b.actions[state] = f
}

// Exec находит для юзера его текущий стейт и исполняет соответствующую функцию
func (b *BotActionsPool) Exec(update tgbotapi.Update) error {
	userState := b.storage.GetState(update.Message.Chat.ID)

	b.actionsMutex.Lock()

	f, ok := b.actions[userState]
	if !ok {
		b.actionsMutex.Unlock()

		return ErrNoAction
	}

	b.actionsMutex.Unlock()

	currentState, err := f(update)
	if err != nil {
		return err
	}

	err = b.storage.SetState(update.Message.Chat.ID, currentState)
	if err != nil {
		return err
	}

	return err
}
