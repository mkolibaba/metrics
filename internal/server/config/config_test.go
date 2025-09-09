package config

import (
	"flag"
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoadServerConfig(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

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

var configFileContent = `{
"address": "some_address",
"crypto_key": "super secret key"
}`

func TestConfigFile(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// создаем временный файл конфигурации
	tmpFile, err := os.CreateTemp("", "test-config-*.json")
	testutils.AssertNoError(t, err)
	defer os.Remove(tmpFile.Name())

	// записываем контент
	_, err = tmpFile.WriteString(configFileContent)
	testutils.AssertNoError(t, err)

	err = tmpFile.Close()
	testutils.AssertNoError(t, err)

	os.Args = append(os.Args, fmt.Sprintf("-config=%s", tmpFile.Name()))

	cfg, err := LoadServerConfig()
	testutils.AssertNoError(t, err)

	if cfg.ServerAddress != "some_address" {
		t.Errorf("want server_address = 'some_address', got %s", cfg.ServerAddress)
	}
	if cfg.CryptoKey != "super secret key" {
		t.Errorf("want crypto_key = 'super secret key', got %s", cfg.CryptoKey)
	}
}
