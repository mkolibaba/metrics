package app

import (
	"context"
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client"
	"github.com/mkolibaba/metrics/internal/agent/sender"
	"github.com/mkolibaba/metrics/internal/common/log"
	stdlog "log"
)

func Run() {
	ctx := context.Background()

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		stdlog.Fatalf("error creating config: %v", err)
	}

	logger := log.New()

	c := collector.NewMetricsCollector(cfg.PollInterval, logger)
	chGauges, chCounters := c.StartCollect(ctx)

	serverAPI := client.New(cfg.ServerAddress, cfg.Key, logger)
	metricsSender := sender.NewMetricsSender(serverAPI, cfg.ReportInterval, cfg.RateLimit, logger)

	logger.Info("running agent")

	metricsSender.StartSend(ctx, chGauges, chCounters)
}
