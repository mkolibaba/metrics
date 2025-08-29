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
	rateLimitDefault             = 100
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY"`
}

type configAlias struct {
	AgentConfig
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

// LoadAgentConfig загружает конфигурацию агента из переменных окружения и флагов командной строки.
// Флаги имеют приоритет над переменными окружения.
func LoadAgentConfig() (*AgentConfig, error) {
	var cfg configAlias

	flag.StringVar(&cfg.ServerAddress, "a", serverAddressDefault, "server address")
	flag.IntVar(&cfg.ReportInterval, "r", reportIntervalSecondsDefault, "report interval (seconds)")
	flag.IntVar(&cfg.PollInterval, "p", pollIntervalSecondsDefault, "poll interval (seconds)")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", rateLimitDefault, "rate limit")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "crypto key file path")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	result := cfg.AgentConfig
	result.ReportInterval = time.Duration(cfg.ReportInterval) * time.Second
	result.PollInterval = time.Duration(cfg.PollInterval) * time.Second

	return &result, nil
}
