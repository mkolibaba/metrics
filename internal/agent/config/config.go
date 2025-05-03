package config

import (
	"flag"
	"github.com/mkolibaba/metrics/internal/common/logger"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func MustLoadAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}
	var reportIntervalString, pollIntervalString string

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&reportIntervalString, "r", "10", "report interval (seconds)")
	flag.StringVar(&pollIntervalString, "p", "2", "poll interval (seconds)")
	flag.Parse()

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = address
	}
	if reportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		reportIntervalString = reportInterval
	}
	if pollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		pollIntervalString = pollInterval
	}

	cfg.ReportInterval = stringToDuration(reportIntervalString)
	cfg.PollInterval = stringToDuration(pollIntervalString)

	return cfg
}

func stringToDuration(s string) time.Duration {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Sugared.Fatalf("error parsing config value: %v", err)
	}
	return time.Duration(i) * time.Second
}
