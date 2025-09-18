package collector

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"go.uber.org/zap/zaptest"
)

func TestMetricsCollector_StartCollect(t *testing.T) {
	pollInterval := 100 * time.Millisecond
	timeout := 150 * time.Millisecond
	maxTimeout := 200 * time.Millisecond

	// Create a collector with a short poll interval for testing
	collector := NewMetricsCollector(pollInterval, zaptest.NewLogger(t).Sugar())

	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), timeout)
		defer cancel()

		// Start collecting in a goroutine
		chGauges, chCounters := collector.StartCollect(ctx)

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
		case <-time.After(maxTimeout):
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
		case <-time.After(maxTimeout):
			t.Error("Timeout waiting for counter metrics")
		}
	})
}

func TestMetricsCollector_ShouldNotCollect(t *testing.T) {
	pollInterval := 100 * time.Millisecond
	timeout := 50 * time.Millisecond
	maxTimeout := 200 * time.Millisecond

	collector := NewMetricsCollector(pollInterval, zaptest.NewLogger(t).Sugar())

	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), timeout)
		defer cancel()

		chGauges, chCounters := collector.StartCollect(ctx)

		select {
		case gauges, ok := <-chGauges:
			if ok {
				t.Error("Expected gauge metrics not to be present", gauges)
			}
		case <-time.After(maxTimeout):
		}

		select {
		case counters, ok := <-chCounters:
			if ok {
				t.Error("Expected counter metrics not to be present", counters)
			}
		case <-time.After(maxTimeout):
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
