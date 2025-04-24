package config

import (
	"flag"
	"log"
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

	address, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.ServerAddress = address
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	}

	reportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		duration, err := stringToDuration(reportInterval)
		if err != nil {
			log.Fatalf("error parsing config value: %v", err)
		}
		cfg.ReportInterval = duration
	} else {
		cfg.ReportInterval = 10 * time.Second
		flag.Func("r", "report interval (seconds)", func(s string) error {
			duration, err := stringToDuration(s)
			if err != nil {
				return err
			}
			cfg.ReportInterval = duration
			return nil
		})
	}

	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		duration, err := stringToDuration(pollInterval)
		if err != nil {
			log.Fatalf("error parsing config value: %v", err)
		}
		cfg.PollInterval = duration
	} else {
		cfg.PollInterval = 2 * time.Second
		flag.Func("p", "poll interval (seconds)", func(s string) error {
			duration, err := stringToDuration(s)
			if err != nil {
				return err
			}
			cfg.PollInterval = duration
			return nil
		})
	}

	flag.Parse()

	return cfg
}

func stringToDuration(s string) (time.Duration, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(i) * time.Second, nil
}
