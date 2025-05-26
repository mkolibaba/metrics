package inmemory

import (
	"testing"
)

func TestShouldUpdateCounter(t *testing.T) {
	t.Run("Should_save_counter_value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateCounter(t.Context(), "my", 12)

		assertCounterHasValue(t, store, "my", int64(12))
	})
	t.Run("Should_save_gauge_value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateGauge(t.Context(), "my", 1.2)

		assertGaugeHasValue(t, store, "my", 1.2)
	})
	t.Run("Should_update_counter_value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateCounter(t.Context(), "my", 12)
		store.UpdateCounter(t.Context(), "my", 12)
		store.UpdateCounter(t.Context(), "my1", 12)

		assertCounterHasValue(t, store, "my", int64(24))
	})
	t.Run("Should_update_gauge_value", func(t *testing.T) {
		store := NewMemStorage()
		store.UpdateGauge(t.Context(), "my", 1.2)
		store.UpdateGauge(t.Context(), "my", 0.8)

		assertGaugeHasValue(t, store, "my", 0.8)
	})
}

func assertCounterHasValue(t *testing.T, store *MemStorage, name string, want int64) {
	t.Helper()
	got, ok := store.counters[name]

	if !ok {
		t.Fatalf("no counter by name = '%s' has been found", name)
	}
	if got != want {
		t.Errorf("want counter['%s'] = %d but got %d", name, want, got)
	}
}

func assertGaugeHasValue(t *testing.T, store *MemStorage, name string, want float64) {
	t.Helper()
	got, ok := store.gauges[name]

	if !ok {
		t.Fatalf("no gauge by name = '%s' has been found", name)
	}
	if got != want {
		t.Errorf("want gauge['%s'] = %f but got %f", name, want, got)
	}
}
