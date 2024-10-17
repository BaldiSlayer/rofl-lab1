package models

// UserState - состояние, в котором сейчас находится пользователь
type UserState int

const (
	// EmptyState - состояние пользователя, когда он в первый раз пришел в бота
	EmptyState UserState = iota
	// StartState - начальное состояние пользователя
	StartState
	// WaitForKBQuestion - состояние ожидания от пользователя ввода запроса к базе знаний
	WaitForKBQuestion
	// WaitForKBResponse - состояние ожидания ответа от backend
	WaitForKBResponse
	TRSState
)
