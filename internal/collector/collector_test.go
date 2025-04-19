package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollect(t *testing.T) {
	t.Run("Should collect metrics", func(t *testing.T) {
		collector := NewMetricsCollector(1 * time.Second)

		collector.collect()

		assert.Equal(t, 1, collector.iterations)
		// рандомные метрики
		assert.Contains(t, collector.Gauges, "Alloc")
		assert.Contains(t, collector.Counters, "PollCount")
	})
	t.Run("Should increment iterations", func(t *testing.T) {
		collector := NewMetricsCollector(1 * time.Second)

		collector.collect()
		collector.collect()
		collector.collect()

		assert.Equal(t, 3, collector.iterations)
	})
}
