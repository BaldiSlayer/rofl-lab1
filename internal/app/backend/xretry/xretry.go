package xretry

import (
	"errors"
	"fmt"
)

type Retrier struct {
	f func() error
}

func Retry(f func() error) *Retrier {
	return &Retrier{f: f}
}

func (r *Retrier) Count(count int) error {
	var err error

	for i := 0; i < count; i++ {
		err1 := r.f()
		if err1 == nil {
			return nil
		}

		err = errors.Join(err, err1)
	}

	return fmt.Errorf(
		"the function did not complete correctly in %d calls with errors: %w",
		count,
		err,
	)
}
