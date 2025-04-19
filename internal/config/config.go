package config

import (
	"flag"
	"time"
)

type AgentConfig struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func LoadAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	reportIntervalSeconds := flag.Int("r", 10, "report interval (seconds)")
	pollIntervalSeconds := flag.Int("p", 2, "poll interval (seconds)")
	flag.Parse()
	cfg.ReportInterval = time.Duration(*reportIntervalSeconds) * time.Second
	cfg.PollInterval = time.Duration(*pollIntervalSeconds) * time.Second
	return cfg
}

type ServerConfig struct {
	ServerAddress string
}

func LoadServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.Parse()
	return cfg
}
