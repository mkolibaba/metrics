package collector

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestMetricsCollector_StartCollect(t *testing.T) {
	// Create a collector with a short poll interval for testing
	collector := NewMetricsCollector(100*time.Millisecond, zap.S())

	// Create channels for gauges and counters
	chGauges := make(chan map[string]float64, 1)
	chCounters := make(chan map[string]int64, 1)

	// Start collecting in a goroutine
	go collector.StartCollect(chGauges, chCounters)

	// Wait for first collection
	time.Sleep(350 * time.Millisecond)

	// Check if we received metrics
	select {
	case gauges := <-chGauges:
		// Verify that we have some basic metrics
		if _, ok := gauges["Alloc"]; !ok {
			t.Error("Expected Alloc metric to be present")
		}
		if _, ok := gauges["RandomValue"]; !ok {
			t.Error("Expected RandomValue metric to be present")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout waiting for gauge metrics")
	}

	select {
	case counters := <-chCounters:
		// Verify that PollCount is present and increasing
		if pollCount, ok := counters["PollCount"]; !ok {
			t.Error("Expected PollCount metric to be present")
		} else if pollCount <= 0 {
			t.Error("Expected PollCount to be greater than 0", pollCount)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Timeout waiting for counter metrics")
	}
}
