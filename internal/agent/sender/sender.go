package sender

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

type ServerAPI interface {
	UpdateCounters(counters map[string]int64) error
	UpdateGauges(gauges map[string]float64) error
}

type MetricsSender struct {
	serverAPI      ServerAPI
	reportInterval time.Duration
	logger         *zap.SugaredLogger
}

func NewMetricsSender(serverAPI ServerAPI, reportInterval time.Duration, logger *zap.SugaredLogger) *MetricsSender {
	return &MetricsSender{serverAPI, reportInterval, logger}
}

func (m *MetricsSender) StartSend(chGauges <-chan map[string]float64, chCounters <-chan map[string]int64) {
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()
	for range ticker.C {
		gauges := <-chGauges
		counters := <-chCounters
		m.logger.Debug("sending metrics to server")
		if err := m.send(gauges, counters); err != nil {
			m.logger.Error(err)
		}
	}
}

func (m *MetricsSender) send(gauges map[string]float64, counters map[string]int64) error {
	err := m.serverAPI.UpdateGauges(gauges)
	if err != nil {
		return fmt.Errorf("error during gauges send: %v", err)
	}

	err = m.serverAPI.UpdateCounters(counters)
	if err != nil {
		return fmt.Errorf("error during counters send: %v", err)
	}

	return nil
}
