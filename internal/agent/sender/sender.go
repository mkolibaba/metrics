package sender

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

type ServerAPI interface {
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
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
	for k, v := range gauges {
		err := m.serverAPI.UpdateGauge(k, v)
		if err != nil {
			return fmt.Errorf("error during gauge value send: %v", err)
		}
	}
	for k, v := range counters {
		err := m.serverAPI.UpdateCounter(k, v)
		if err != nil {
			return fmt.Errorf("error during counter value send: %v", err)
		}
	}
	return nil
}
