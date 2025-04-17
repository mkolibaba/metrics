package sender

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/collector"
	"log"
	"net/http"
	"strconv"
	"time"
)

const reportInterval = 10 * time.Second

type MetricsSender struct {
	collector *collector.MetricsCollector
}

func NewMetricsSender(collector *collector.MetricsCollector) *MetricsSender {
	return &MetricsSender{collector}
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
		sendMetric("gauge", k, strconv.FormatFloat(v, 'f', 4, 64))
	}
	for k, v := range m.collector.Counters {
		sendMetric("counter", k, strconv.FormatInt(v, 10))
	}
}

func sendMetric(t, name, val string) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", t, name, val)
	fmt.Printf("Sending POST %s\n", url)
	_, err := http.Post(url, "text/plain", nil)
	if err != nil {
		log.Print(err)
	}
}
