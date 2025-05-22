package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	ServerAddress   string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
}

func MustLoadServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	var storeIntervalSeconds int

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&storeIntervalSeconds, "i", 300, "store interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres:postgres@localhost:5432/metrics", "server address")
	flag.Parse()

	if address, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = address
	}
	if v, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeIntervalSeconds = mustStringToInt(v)
	}
	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = v
	}
	if v, ok := os.LookupEnv("RESTORE"); ok {
		cfg.Restore = mustStringToBool(v)
	}
	if v, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = v
	}

	cfg.StoreInterval = time.Duration(storeIntervalSeconds) * time.Second

	return cfg
}

func mustStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("error parsing config value: %v", err)
	}
	return i
}

func mustStringToBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("error parsing config value: %v", err)
	}
	return b
}
