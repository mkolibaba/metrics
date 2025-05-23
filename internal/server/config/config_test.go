package config

import (
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	address := "123"
	t.Setenv("ADDRESS", address)

	storeInterval := "2"
	storeIntervalDuration := 2 * time.Second
	t.Setenv("STORE_INTERVAL", storeInterval)

	fileStoragePath := "/a/b/c/file.txt"
	t.Setenv("FILE_STORAGE_PATH", fileStoragePath)

	restore := "true"
	restoreBool := true
	t.Setenv("RESTORE", restore)

	dbDSN := "my-db-dsn"
	t.Setenv("DATABASE_DSN", dbDSN)

	cfg, err := LoadServerConfig()
	testutils.AssertNoError(t, err)

	if cfg.ServerAddress != address {
		t.Errorf("want server address %s, got %s", address, cfg.ServerAddress)
	}
	if cfg.StoreInterval != storeIntervalDuration {
		t.Errorf("want store interval %s, got %s", storeIntervalDuration, cfg.StoreInterval)
	}
	if cfg.FileStoragePath != fileStoragePath {
		t.Errorf("want file storage path %s, got %s", fileStoragePath, cfg.FileStoragePath)
	}
	if cfg.Restore != restoreBool {
		t.Errorf("want restore %v, got %v", restoreBool, cfg.Restore)
	}
	if cfg.DatabaseDSN != dbDSN {
		t.Errorf("want database dsn %s, got %s", dbDSN, cfg.DatabaseDSN)
	}
}
