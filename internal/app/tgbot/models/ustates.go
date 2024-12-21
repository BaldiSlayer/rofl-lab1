package models

// UserState - состояние, в котором сейчас находится пользователь
type UserState int

const (
	// Start - состояние пользователя, когда он в первый раз пришел в бота
	Start UserState = iota
	// GetRequest - состояние ожидания ввода запроса от пользователя
	GetRequest
	// GetTrs - сосояние ожидания ввода TRS от пользователя
	GetTrs
	// ValidateTrs - сосотояние ожидания подтверждения корректности выделения
	// формальной TRS
	ValidateTrs
	// FixTrs - состояние отображения пользователю ошибки формализации TRS
	FixTrs
	// GetQuestionMultiModels -
	GetQuestionMultiModels
)
