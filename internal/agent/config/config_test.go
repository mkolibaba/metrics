package config

import (
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	address := "123123"
	t.Setenv("ADDRESS", address)

	reportInterval := "33"
	reportIntervalDuration := 33 * time.Second
	t.Setenv("REPORT_INTERVAL", reportInterval)

	pollInterval := "22"
	pollIntervalDuration := 22 * time.Second
	t.Setenv("POLL_INTERVAL", pollInterval)

	cfg, err := LoadAgentConfig()
	testutils.AssertNoError(t, err)

	if cfg.ServerAddress != address {
		t.Errorf("want server address %s, got %s", address, cfg.ServerAddress)
	}
	if cfg.ReportInterval != reportIntervalDuration {
		t.Errorf("want report interval %s, got %s", reportIntervalDuration, cfg.ReportInterval)
	}
	if cfg.PollInterval != pollIntervalDuration {
		t.Errorf("want poll interval %s, got %s", pollIntervalDuration, cfg.PollInterval)
	}
}
