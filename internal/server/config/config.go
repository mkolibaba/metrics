package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"time"
)

type ServerConfig struct {
	ServerAddress   string `env:"ADDRESS"`
	StoreInterval   time.Duration
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}
type configAlias struct {
	ServerConfig
	StoreInterval int `env:"STORE_INTERVAL"`
}

func LoadServerConfig() (*ServerConfig, error) {
	var cfg configAlias

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "server address")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	result := cfg.ServerConfig
	result.StoreInterval = time.Duration(cfg.StoreInterval) * time.Second

	return &result, nil
}
