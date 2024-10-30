package models

// UserState - состояние, в котором сейчас находится пользователь
type UserState int

const (
	// EmptyState - состояние пользователя, когда он в первый раз пришел в бота
	EmptyState UserState = iota
	// WaitForRequest - состояние ожидания ввода запроса от пользователя
	WaitForRequest
	// WaitForKBResponse - состояние ожидания ответа от backend
	WaitForKBResponse
	TRSState
)
