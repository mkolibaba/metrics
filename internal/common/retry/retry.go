package retry

import "time"

const retryAttempts = 3

var retryIntervalsSeconds = []int{1, 3, 5}

func Do(fn func() error) error {
	var err error
	for i := 0; i <= retryAttempts; i++ {
		err = fn()
		if err == nil || i == retryAttempts {
			break
		}

		interval := time.Duration(retryIntervalsSeconds[i]) * time.Second
		time.Sleep(interval)
	}

	return err
}

func DoWithReturn[T any](fn func() (T, error)) (T, error) {
	var (
		value T
		err   error
	)
	for i := 0; i <= retryAttempts; i++ {
		value, err = fn()
		if err == nil || i == retryAttempts {
			break
		}

		interval := time.Duration(retryIntervalsSeconds[i]) * time.Second
		time.Sleep(interval)
	}

	return value, err
}
