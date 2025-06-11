package sender

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

//go:generate moq -stub -out server_api_mock.go . ServerAPI
type ServerAPI interface {
	UpdateCounters(counters map[string]int64) error
	UpdateGauges(gauges map[string]float64) error
}

type MetricsSender struct {
	serverAPI       ServerAPI
	reportInterval  time.Duration
	reportRateLimit int
	logger          *zap.SugaredLogger
}

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

loop:
	for {
		select {
		case <-ticker.C:
			m.logger.Debug("starting to send metrics to server")

			aggregated := m.aggregate(chGauges, chCounters)
			for i := 0; i < m.reportRateLimit; i++ {
				go worker(aggregated)
			}
		case <-ctx.Done():
			break loop
		}
	}
}

func (m *MetricsSender) aggregate(
	chGauges <-chan map[string]float64,
	chCounters <-chan map[string]int64,
) <-chan func() error {
	c := make(chan func() error)

	go func() {
	loop:
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
			// все значение вычитаны, агрегирующий канал можно закрывать
			default:
				m.logger.Debug("closing aggregating channel")
				close(c)
				break loop
			}
		}
	}()

	return c
}
