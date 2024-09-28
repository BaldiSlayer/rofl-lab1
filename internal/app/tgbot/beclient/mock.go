package beclient

import (
	"context"
	"time"
)

type MockBackendClient struct{}

func (mbc *MockBackendClient) AskKB(ctx context.Context, question string) (string, error) {
	time.Sleep(1 * time.Second)

	return "пошел нахуй", nil
}

func (mbc *MockBackendClient) ParseTRS(ctx context.Context, trs string) error {
	return nil
}
