package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	ServerAddress string
}

func MustLoadServerConfig() *ServerConfig {
	cfg := &ServerConfig{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.Parse()

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = address
	}

	return cfg
}
