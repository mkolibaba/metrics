package collector

import (
	"testing"
	"testing/synctest"
	"time"

	"go.uber.org/zap/zaptest"
)

func TestMetricsCollector_StartCollect(t *testing.T) {
	pollInterval := 100 * time.Millisecond
	waitDuration := 150 * time.Millisecond

	// Create a collector with a short poll interval for testing
	collector := NewMetricsCollector(pollInterval, zaptest.NewLogger(t).Sugar())

	synctest.Test(t, func(t *testing.T) {
		// Start collecting in a goroutine
		chGauges, chCounters := collector.StartCollect(t.Context())

		time.Sleep(waitDuration)

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
		default:
			t.Error("Expected gauges to be present")
		}

		select {
		case counters := <-chCounters:
			// Verify that PollCount is present and increasing
			if pollCount, ok := counters["PollCount"]; !ok {
				t.Error("Expected PollCount metric to be present")
			} else if pollCount <= 0 {
				t.Error("Expected PollCount to be greater than 0", pollCount)
			}
		default:
			t.Error("Expected counters to be present")
		}
	})
}

func TestMetricsCollector_ShouldNotCollect(t *testing.T) {
	pollInterval := 100 * time.Millisecond
	waitDuration := 50 * time.Millisecond

	collector := NewMetricsCollector(pollInterval, zaptest.NewLogger(t).Sugar())

	synctest.Test(t, func(t *testing.T) {
		chGauges, chCounters := collector.StartCollect(t.Context())

		time.Sleep(waitDuration)

		gaugesArePresent := false

		select {
		case _, ok := <-chGauges:
			gaugesArePresent = ok
		default:
		}

		if gaugesArePresent {
			t.Error("Expected gauge metrics not to be present")
		}

		countersArePresent := false

		select {
		case _, ok := <-chCounters:
			countersArePresent = ok
		default:
		}

		if countersArePresent {
			t.Error("Expected counter metrics not to be present")
		}
	})
}

func BenchmarkMetricsCollector_CollectAdditionalGauges(b *testing.B) {
	collector := NewMetricsCollector(100*time.Millisecond, zaptest.NewLogger(b).Sugar())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collector.collectAdditionalGauges()
	}
}
