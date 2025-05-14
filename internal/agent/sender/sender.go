package sender

import (
	"fmt"
	"time"

	"github.com/mkolibaba/metrics/internal/common/logger"
)

type ServerAPI interface {
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
}

type MetricsSender struct {
	serverAPI      ServerAPI
	reportInterval time.Duration
}

func NewMetricsSender(serverAPI ServerAPI, reportInterval time.Duration) *MetricsSender {
	return &MetricsSender{serverAPI, reportInterval}
}

func (m *MetricsSender) StartSend(chGauges <-chan map[string]float64, chCounters <-chan map[string]int64) error {
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()
	for range ticker.C {
		gauges := <-chGauges
		counters := <-chCounters
		if err := m.send(gauges, counters); err != nil {
			logger.Sugared.Errorf("error during metrics send: %v", err)
			return err
		}
	}
	return nil
}

func (m *MetricsSender) send(gauges map[string]float64, counters map[string]int64) error {
	for k, v := range gauges {
		err := m.serverAPI.UpdateGauge(k, v)
		if err != nil {
			logger.Sugared.Errorf("error during gauge value send: %v", err)
			return fmt.Errorf("error during gauge value send: %v", err)
		}
	}
	for k, v := range counters {
		err := m.serverAPI.UpdateCounter(k, v)
		if err != nil {
			logger.Sugared.Errorf("error during counter value send: %v", err)
			return fmt.Errorf("error during counter value send: %v", err)
		}
	}
	return nil
}
