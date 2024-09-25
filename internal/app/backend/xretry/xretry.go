package xretry

import (
	"errors"
	"fmt"
	"time"
)

type Retrier struct {
	f     func() error
	delay time.Duration
}

func Retry(f func() error) *Retrier {
	return &Retrier{f: f}
}

func (r *Retrier) WithDelay(delay time.Duration) *Retrier {
	r.delay = delay

	return r
}

func (r *Retrier) Count(count int) error {
	var err error

	for i := 0; i < count; i++ {
		err1 := r.f()
		if err1 == nil {
			return nil
		}

		err = errors.Join(err, err1)

		// чтобы не ждать после последнего выполнения
		if i == count-1 {
			break
		}

		time.Sleep(r.delay)
	}

	return fmt.Errorf(
		"the function did not complete correctly in %d calls with errors: %w",
		count,
		err,
	)
}
