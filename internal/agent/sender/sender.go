// Package sender предоставляет компоненты для периодической отправки
// собранных метрик на сервер.
package sender

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ServerAPI описывает клиентский интерфейс для отправки метрик на сервер.
// Реализация должна уметь передавать счётчики и gauge метрики.
//
//go:generate moq -stub -out server_api_mock.go . ServerAPI
type ServerAPI interface {
	UpdateCounters(counters map[string]int64) error
	UpdateGauges(gauges map[string]float64) error
}

// MetricsSender выполняет периодическую отправку метрик на сервер.
type MetricsSender struct {
	serverAPI       ServerAPI
	reportInterval  time.Duration
	reportRateLimit int
	logger          *zap.SugaredLogger
}

// NewMetricsSender создаёт и настраивает отправитель метрик.
// reportInterval — период отправки, reportRateLimit — число параллельных воркеров для отправки.
func NewMetricsSender(
	serverAPI ServerAPI,
	reportInterval time.Duration,
	reportRateLimit int,
	logger *zap.SugaredLogger,
) *MetricsSender {
	return &MetricsSender{
		serverAPI:       serverAPI,
		reportInterval:  reportInterval,
		reportRateLimit: reportRateLimit,
		logger:          logger,
	}
}

// StartSend запускает цикл отправки метрик на сервер. Метод принимает
// каналы gauge метрик (map[string]float64) и counter метрик
// (map[string]int64), агрегирует задачи отправки и выполняет их
// указанным числом воркеров. Останавливается при отмене контекста.
func (m *MetricsSender) StartSend(
	ctx context.Context,
	chGauges <-chan map[string]float64,
	chCounters <-chan map[string]int64,
) {
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()

	worker := func(in <-chan func() error) {
		for fn := range in {
			if err := fn(); err != nil {
				m.logger.Errorf("failed to execute sending: %s", err)
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			m.logger.Debug("starting to send metrics to server")

			aggregated := m.aggregate(ctx, chGauges, chCounters)
			for i := 0; i < m.reportRateLimit; i++ {
				go worker(aggregated)
			}
		case <-ctx.Done():
			m.logger.Debug("stopping sending")
			return
		}
	}
}

func (m *MetricsSender) aggregate(
	ctx context.Context,
	chGauges <-chan map[string]float64,
	chCounters <-chan map[string]int64,
) <-chan func() error {
	c := make(chan func() error)
	closeC := func() {
		m.logger.Debug("closing aggregating channel")
		close(c)
	}

	go func() {
		for {
			select {
			// вычитываем все из канала gauges
			case v := <-chGauges:
				m.logger.Debug("reading from gauges channel")
				c <- func() error {
					err := m.serverAPI.UpdateGauges(v)
					if err != nil {
						return fmt.Errorf("send gauges: %w", err)
					}
					return nil
				}
			// вычитываем все из канала counters
			case v := <-chCounters:
				m.logger.Debug("reading from counters channel")
				c <- func() error {
					err := m.serverAPI.UpdateCounters(v)
					if err != nil {
						return fmt.Errorf("send counters: %w", err)
					}
					return nil
				}
			case <-ctx.Done():
				closeC()
				return
			// все значение вычитаны, агрегирующий канал можно закрывать
			default:
				closeC()
				return
			}
		}
	}()

	return c
}
