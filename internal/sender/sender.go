package sender

import (
	"github.com/mkolibaba/metrics/internal/collector"
	"github.com/mkolibaba/metrics/internal/http/client"
	"log"
	"time"
)

const reportInterval = 10 * time.Second

type MetricsSender struct {
	collector *collector.MetricsCollector
	serverAPI client.ServerAPI
}

func NewMetricsSender(collector *collector.MetricsCollector, serverAPI client.ServerAPI) *MetricsSender {
	return &MetricsSender{collector, serverAPI}
}

func (m *MetricsSender) StartCollectAndSend() {
	m.collector.StartCollect()
	go func() {
		for {
			time.Sleep(reportInterval)
			m.send()
		}
	}()
}

func (m *MetricsSender) send() {
	for k, v := range m.collector.Gauges {
		err := m.serverAPI.UpdateGauge(k, v)
		if err != nil {
			log.Print(err)
		}
	}
	for k, v := range m.collector.Counters {
		err := m.serverAPI.UpdateCounter(k, v)
		if err != nil {
			log.Print(err)
		}
	}
}
