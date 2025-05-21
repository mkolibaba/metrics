package sender

import (
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/mkolibaba/metrics/internal/agent/http/client/mocks"
)

func TestStartSend(t *testing.T) {
	// Create mock server API
	serverAPI := &mocks.ServerAPIMock{}

	// Create channels for metrics
	chGauges := make(chan map[string]float64, 1)
	chCounters := make(chan map[string]int64, 1)

	// Create sender with short report interval for testing
	sender := NewMetricsSender(serverAPI, 100*time.Millisecond, zap.S())

	// Start sending in a goroutine
	go func() {
		sender.StartSend(chGauges, chCounters)
	}()

	// Send test metrics
	testGauges := map[string]float64{"test_gauge": 42.0}
	testCounters := map[string]int64{"test_counter": 1}

	// Send metrics through channels
	chGauges <- testGauges
	chCounters <- testCounters

	// Wait for a short time to allow processing
	time.Sleep(150 * time.Millisecond)

	// Verify that the metrics were sent
	if serverAPI.GaugeCalls != 1 {
		t.Errorf("expected 1 gauge call, got %d", serverAPI.GaugeCalls)
	}
	if serverAPI.CounterCalls != 1 {
		t.Errorf("expected 1 counter call, got %d", serverAPI.CounterCalls)
	}
}
