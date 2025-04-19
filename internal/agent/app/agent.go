package app

import (
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client"
	"github.com/mkolibaba/metrics/internal/agent/sender"
	"log"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.LoadAgentConfig()

	c := collector.NewMetricsCollector(cfg.PollInterval)
	serverAPI := client.New(cfg.ServerAddress)

	metricsSender := sender.NewMetricsSender(c, serverAPI, cfg.ReportInterval)
	metricsSender.StartCollectAndSend()

	select {}
}
