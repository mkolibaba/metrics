package config

import (
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"reflect"
	"testing"
	"time"
)

func TestLoadServerConfig(t *testing.T) {
	t.Setenv("ADDRESS", "123")
	t.Setenv("STORE_INTERVAL", "2")
	t.Setenv("FILE_STORAGE_PATH", "/a/b/c/file.txt")
	t.Setenv("RESTORE", "true")
	t.Setenv("DATABASE_DSN", "my-db-dsn")
	wantConfig := &ServerConfig{
		ServerAddress:   "123",
		StoreInterval:   2 * time.Second,
		FileStoragePath: "/a/b/c/file.txt",
		Restore:         true,
		DatabaseDSN:     "my-db-dsn",
	}

	gotConfig, err := LoadServerConfig()
	testutils.AssertNoError(t, err)

	if !reflect.DeepEqual(wantConfig, gotConfig) {
		t.Errorf("want config %#v, got %#v", wantConfig, gotConfig)
	}
}
