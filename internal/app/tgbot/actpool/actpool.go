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

func New(storage ustorage.UserDataStorage) (*BotActionsPool, error) {
	return &BotActionsPool{
		actions: make(map[models.UserState]StateTransition),
		storage: storage,
	}, nil
}

// AddStateTransition добавляет контроллер. Не использовать "на лету", добавлять контроллеры только
// при запуске программы. Иначе получите Race Condition. Не добавляю mutex, потому что добавлять
// контроллеры "на лету" - очень странный кейс
// FIXME: передавать сразу мапу в конструктор, вместо этой функции
func (b *BotActionsPool) AddStateTransition(state models.UserState, f StateTransition) {
	b.actions[state] = f
}

// Exec находит для юзера его текущий стейт и исполняет соответствующую функцию
func (b *BotActionsPool) Exec(update tgbotapi.Update) error {
	// FIXME: брать мьютекс на юзера

	userState, err := b.storage.GetState(update.Message.Chat.ID)
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

	err = b.storage.SetState(update.Message.Chat.ID, currentState)
	if err != nil {
		return err
	}

	return err
}
