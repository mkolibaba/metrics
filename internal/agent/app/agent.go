package app

import (
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client"
	"github.com/mkolibaba/metrics/internal/agent/sender"
	"github.com/mkolibaba/metrics/internal/common/logger"
)

func Run() {
	cfg := config.MustLoadAgentConfig()

	chGauges := make(chan map[string]float64, 1)
	chCounters := make(chan map[string]int64, 1)

	c := collector.NewMetricsCollector(cfg.PollInterval)
	go c.StartCollect(chGauges, chCounters)

	serverAPI := client.New(cfg.ServerAddress)
	metricsSender := sender.NewMetricsSender(serverAPI, cfg.ReportInterval)
	if err := metricsSender.StartSend(chGauges, chCounters); err != nil {
		logger.Sugared.Errorf("error during metrics send: %v", err)
	}
}
