package app

import (
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client"
	"github.com/mkolibaba/metrics/internal/agent/sender"
	"github.com/mkolibaba/metrics/internal/common/log"
)

func Run() {
	cfg := config.MustLoadAgentConfig()

	logger := log.New()

	chGauges := make(chan map[string]float64, 1)
	chCounters := make(chan map[string]int64, 1)

	c := collector.NewMetricsCollector(cfg.PollInterval, logger)
	go c.StartCollect(chGauges, chCounters)

	serverAPI := client.New(cfg.ServerAddress)
	metricsSender := sender.NewMetricsSender(serverAPI, cfg.ReportInterval, logger)

	logger.Info("running agent")

	metricsSender.StartSend(chGauges, chCounters)
}
