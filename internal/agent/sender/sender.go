package sender

import (
	"log"
	"time"
)

type Collector interface {
	StartCollect()
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}

type ServerAPI interface {
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
}

type MetricsSender struct {
	collector      Collector
	serverAPI      ServerAPI
	reportInterval time.Duration
}

func NewMetricsSender(collector Collector, serverAPI ServerAPI, reportInterval time.Duration) *MetricsSender {
	return &MetricsSender{collector, serverAPI, reportInterval}
}

func (m *MetricsSender) StartCollectAndSend() {
	m.collector.StartCollect()
	m.StartSend()
}

func (m *MetricsSender) StartSend() {
	go func() {
		for {
			time.Sleep(m.reportInterval)
			m.send()
		}
	}()
}

func (m *MetricsSender) send() {
	for k, v := range m.collector.GetGauges() {
		err := m.serverAPI.UpdateGauge(k, v)
		if err != nil {
			log.Print(err)
		}
	}
	for k, v := range m.collector.GetCounters() {
		err := m.serverAPI.UpdateCounter(k, v)
		if err != nil {
			log.Print(err)
		}
	}
}
