package collector

import (
	"testing"
	"time"
)

func TestCollect(t *testing.T) {
	t.Run("Should_collect_metrics", func(t *testing.T) {
		collector := NewMetricsCollector(1 * time.Second)

		collector.collect()

		assertCollectorIterations(t, collector, 1)

		// рандомные метрики
		wantGauge := "Alloc"
		if _, ok := collector.gauges[wantGauge]; !ok {
			t.Errorf("no gauge by name = '%s' has been found", wantGauge)
		}
		wantCounter := "PollCount"
		if _, ok := collector.counters[wantCounter]; !ok {
			t.Errorf("no counter by name = '%s' has been found", wantCounter)
		}
	})
	t.Run("Should_increment_iterations", func(t *testing.T) {
		collector := NewMetricsCollector(1 * time.Second)

		collector.collect()
		collector.collect()
		collector.collect()

		assertCollectorIterations(t, collector, 3)
	})
}

func assertCollectorIterations(t *testing.T, collector *MetricsCollector, want int) {
	t.Helper()
	got := collector.iterations
	if got != want {
		t.Errorf("want collector to record %d iterations, got %d", want, got)
	}
}
