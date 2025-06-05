package app

import (
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client"
	"github.com/mkolibaba/metrics/internal/agent/sender"
	"github.com/mkolibaba/metrics/internal/common/log"
	stdlog "log"
)

func Run() {
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		stdlog.Fatalf("error creating config: %v", err)
	}

	logger := log.New()

	chGauges := make(chan map[string]float64, 1)
	chCounters := make(chan map[string]int64, 1)

	c := collector.NewMetricsCollector(cfg.PollInterval, logger)
	go c.StartCollect(chGauges, chCounters)

	serverAPI := client.New(cfg.ServerAddress, logger)
	metricsSender := sender.NewMetricsSender(serverAPI, cfg.ReportInterval, logger)

	logger.Info("running agent")

	metricsSender.StartSend(chGauges, chCounters)
}
