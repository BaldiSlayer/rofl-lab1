package mclient

type ModelClient interface {
	// Ask отправляет запрос к модели
	Ask(question string) (string, error)
	// AskWithContext отправляет запрос к модели с использованием контекста
	AskWithContext(question string) (string, error)
}
