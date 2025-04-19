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
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "poll interval")
	flag.Parse()
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
