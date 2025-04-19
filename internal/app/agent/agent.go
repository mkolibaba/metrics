package agent

import (
	"github.com/mkolibaba/metrics/internal/collector"
	"github.com/mkolibaba/metrics/internal/config"
	"github.com/mkolibaba/metrics/internal/http/client"
	"github.com/mkolibaba/metrics/internal/sender"
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
