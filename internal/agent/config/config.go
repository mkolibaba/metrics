package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"time"
)

const (
	serverAddressDefault         = "localhost:8080"
	reportIntervalSecondsDefault = 10
	pollIntervalSecondsDefault   = 10
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string `env:"KEY"`
}

type configAlias struct {
	AgentConfig
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

func LoadAgentConfig() (*AgentConfig, error) {
	var cfg configAlias

	flag.StringVar(&cfg.ServerAddress, "a", serverAddressDefault, "server address")
	flag.IntVar(&cfg.ReportInterval, "r", reportIntervalSecondsDefault, "report interval (seconds)")
	flag.IntVar(&cfg.PollInterval, "p", pollIntervalSecondsDefault, "poll interval (seconds)")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	result := cfg.AgentConfig
	result.ReportInterval = time.Duration(cfg.ReportInterval) * time.Second
	result.PollInterval = time.Duration(cfg.PollInterval) * time.Second

	return &result, nil
}
