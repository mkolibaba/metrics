// Package retry предоставляет функционал повторного выполнения операций
// с фиксированными интервалами ожидания между попытками.
package retry

import "time"

const retryAttempts = 3

var retryIntervalsSeconds = []int{1, 3, 5}

// Do выполняет fn с несколькими повторными попытками при ошибке.
// Возвращает последнюю ошибку либо nil, если одна из попыток успешна.
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

// DoWithReturn выполняет fn и повторяет попытки при ошибке, возвращая
// значение и ошибку последней попытки (или успешное значение и nil).
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
