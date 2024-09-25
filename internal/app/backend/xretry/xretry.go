package xretry

import "errors"

type Retrier struct {
	f func() error
}

func Retry(f func() error) *Retrier {
	return &Retrier{f: f}
}

func (r *Retrier) Count(count int) error {
	for i := 0; i < count; i++ {
		err := r.f()
		if err == nil {
			return nil
		}
	}

	// не получилось за count попыток
	return errors.New("")
}
