package main

import (
	"github.com/mkolibaba/metrics/internal/collector"
	"github.com/mkolibaba/metrics/internal/sender"
)

func main() {
	metricsSender := sender.NewMetricsSender(collector.NewMetricsCollector())
	metricsSender.StartCollectAndSend()

	select {}
}
