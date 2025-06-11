package sender

import (
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

func TestStartSend(t *testing.T) {
	// Create mock server API
	serverAPI := &ServerAPIMock{}

	// Create channels for metrics
	chGauges := make(chan map[string]float64, 1)
	chCounters := make(chan map[string]int64, 1)

	// Create sender with short report interval for testing
	sender := NewMetricsSender(serverAPI, 100*time.Millisecond, 10, zaptest.NewLogger(t).Sugar())

	// Start sending in a goroutine
	go func() {
		sender.StartSend(t.Context(), chGauges, chCounters)
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
	if len(serverAPI.UpdateGaugesCalls()) != 1 {
		t.Errorf("expected 1 gauge call, got %d", len(serverAPI.UpdateGaugesCalls()))
	}
	if len(serverAPI.UpdateCountersCalls()) != 1 {
		t.Errorf("expected 1 counter call, got %d", len(serverAPI.UpdateCountersCalls()))
	}
}
