package beclient

import "context"

type BackendClient interface {
	// AskKB отправляет запрос для получения ответа с помощью базы знаний
	AskKB(ctx context.Context, question string) (string, error)
	// ParseTRS отправляет запрос для TRS
	ParseTRS(ctx context.Context, trs string) error
}
