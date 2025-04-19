package config

import (
	"flag"
	"os"
)

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
