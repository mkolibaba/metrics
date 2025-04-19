package agent

import (
	"github.com/mkolibaba/metrics/internal/collector"
	"github.com/mkolibaba/metrics/internal/http/client"
	"github.com/mkolibaba/metrics/internal/sender"
)

func Run() {
	metricsSender := sender.NewMetricsSender(collector.NewMetricsCollector(), &client.ServerClient{})
	metricsSender.StartCollectAndSend()

	select {}
}
