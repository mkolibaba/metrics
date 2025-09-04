package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v11"
	"os"
	"strings"
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

type rawConfig struct {
	AgentConfig
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

// LoadAgentConfig загружает конфигурацию агента. Значения имеют следующий приоритет:
// переменные окружения > флаги > значения из конфигурационного файла > значения по умолчанию.
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := createDefaultConfig()

	if err := readFromConfigFile(&cfg); err != nil {
		return nil, err
	}

	parseFlags(&cfg)

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	result := cfg.AgentConfig
	result.ReportInterval = time.Duration(cfg.ReportInterval) * time.Second
	result.PollInterval = time.Duration(cfg.PollInterval) * time.Second

	return &result, nil
}

func createDefaultConfig() rawConfig {
	var cfg rawConfig
	cfg.ServerAddress = serverAddressDefault
	cfg.ReportInterval = reportIntervalSecondsDefault
	cfg.PollInterval = pollIntervalSecondsDefault
	cfg.RateLimit = rateLimitDefault
	return cfg
}

func readFromConfigFile(cfg *rawConfig) error {
	fn := func() string {
		if configFilePath, ok := os.LookupEnv("CONFIG"); ok {
			return configFilePath
		}

		args := os.Args[1:]
		for i, arg := range args {
			if arg == "-c" || arg == "-config" {
				return args[i+1]
			}
			if strings.HasPrefix(arg, "-c=") {
				return strings.TrimPrefix(arg, "-c=")
			}
			if strings.HasPrefix(arg, "-config=") {
				return strings.TrimPrefix(arg, "-config=")
			}
		}

		return ""
	}

	path := fn()
	if path == "" {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, cfg)
}

func parseFlags(cfg *rawConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", serverAddressDefault, "server address")
	flag.IntVar(&cfg.ReportInterval, "r", reportIntervalSecondsDefault, "report interval (seconds)")
	flag.IntVar(&cfg.PollInterval, "p", pollIntervalSecondsDefault, "poll interval (seconds)")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", rateLimitDefault, "rate limit")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "crypto key file path")
	_ = flag.String("c", "", "config file path")
	_ = flag.String("config", "", "config file path")
	flag.Parse()
}
