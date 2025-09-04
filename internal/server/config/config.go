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
	serverAddressDefault        = "localhost:8080"
	storeIntervalSecondsDefault = 300
	fileStoragePathDefault      = "db.json"
	restoreDefault              = true
)

type ServerConfig struct {
	ServerAddress   string `env:"ADDRESS" json:"address"`
	StoreInterval   time.Duration
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	Key             string `env:"KEY" json:"key"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
}

type rawConfig struct {
	ServerConfig
	StoreInterval int `env:"STORE_INTERVAL" json:"store_interval"`
}

// LoadServerConfig загружает конфигурацию сервера. Значения имеют следующий приоритет:
// переменные окружения > флаги > значения из конфигурационного файла > значения по умолчанию.
func LoadServerConfig() (*ServerConfig, error) {
	cfg := createDefaultConfig()

	if err := readFromConfigFile(&cfg); err != nil {
		return nil, err
	}

	parseFlags(&cfg)

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	result := cfg.ServerConfig
	result.StoreInterval = time.Duration(cfg.StoreInterval) * time.Second

	return &result, nil
}

func createDefaultConfig() rawConfig {
	var cfg rawConfig
	cfg.ServerAddress = serverAddressDefault
	cfg.StoreInterval = storeIntervalSecondsDefault
	cfg.FileStoragePath = fileStoragePathDefault
	cfg.Restore = restoreDefault
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
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.IntVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "store interval")
	flag.StringVar(&cfg.FileStoragePath, cfg.FileStoragePath, "db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "server address")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "hash key")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", cfg.CryptoKey, "crypto key file path")
	_ = flag.String("c", "", "config file path")
	_ = flag.String("config", "", "config file path")
	flag.Parse()
}
