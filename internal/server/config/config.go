package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"time"
)

const (
	serverAddressDefault        = "localhost:8080"
	storeIntervalSecondsDefault = 300
	fileStoragePathDefault      = "db.json"
	restoreDefault              = true
)

type ServerConfig struct {
	ServerAddress   string `env:"ADDRESS"`
	StoreInterval   time.Duration
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}
type configAlias struct {
	ServerConfig
	StoreInterval int `env:"STORE_INTERVAL"`
}

func LoadServerConfig() (*ServerConfig, error) {
	var cfg configAlias

	flag.StringVar(&cfg.ServerAddress, "a", serverAddressDefault, "server address")
	flag.IntVar(&cfg.StoreInterval, "i", storeIntervalSecondsDefault, "store interval")
	flag.StringVar(&cfg.FileStoragePath, fileStoragePathDefault, "db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", restoreDefault, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "server address")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	result := cfg.ServerConfig
	result.StoreInterval = time.Duration(cfg.StoreInterval) * time.Second

	return &result, nil
}
