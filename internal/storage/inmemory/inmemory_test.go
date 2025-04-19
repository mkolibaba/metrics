package inmemory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldUpdateCounter(t *testing.T) {
	t.Run("Should save counter value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateCounter("my", 12)

		assert.Contains(t, store.counters, "my")
		assert.Equal(t, int64(12), store.counters["my"])
	})
	t.Run("Should save gauge value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateGauge("my", 1.2)

		assert.Contains(t, store.gauges, "my")
		assert.Equal(t, 1.2, store.gauges["my"])
	})
	t.Run("Should update counter value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateCounter("my", 12)
		store.UpdateCounter("my", 12)
		store.UpdateCounter("my1", 12)

		assert.Equal(t, int64(24), store.counters["my"])
	})
	t.Run("Should update gauge value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateGauge("my", 1.2)
		store.UpdateGauge("my", 0.8)

		assert.Equal(t, 0.8, store.gauges["my"])
	})
}
