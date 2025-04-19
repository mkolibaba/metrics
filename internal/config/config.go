package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func LoadAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}

	address, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.ServerAddress = address
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	}

	reportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		i, err := strconv.Atoi(reportInterval)
		if err != nil {
			panic(err)
		}
		cfg.ReportInterval = time.Duration(i) * time.Second
	} else {
		cfg.ReportInterval = 10 * time.Second
		flag.Func("r", "report interval (seconds)", func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			cfg.ReportInterval = time.Duration(i) * time.Second
			return nil
		})
	}

	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		i, err := strconv.Atoi(pollInterval)
		if err != nil {
			panic(err)
		}
		cfg.PollInterval = time.Duration(i) * time.Second
	} else {
		cfg.PollInterval = 2 * time.Second
		flag.Func("p", "poll interval (seconds)", func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			cfg.PollInterval = time.Duration(i) * time.Second
			return nil
		})
	}

	flag.Parse()

	return cfg
}

type ServerConfig struct {
	ServerAddress string
}

func LoadServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	address, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.ServerAddress = address
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	}
	flag.Parse()
	return cfg
}
