package update

import (
	"slices"
	"testing"
)

type MetricsUpdaterMock struct {
	Calls                int
	NamesPassed          []string
	GaugesValuesPassed   []float64
	CountersValuesPassed []int64
}

func (m *MetricsUpdaterMock) UpdateGauge(name string, value float64) (float64, error) {
	m.Calls++
	m.NamesPassed = append(m.NamesPassed, name)
	m.GaugesValuesPassed = append(m.GaugesValuesPassed, value)
	return 0, nil
}

func (m *MetricsUpdaterMock) UpdateCounter(name string, value int64) (int64, error) {
	m.Calls++
	m.NamesPassed = append(m.NamesPassed, name)
	m.CountersValuesPassed = append(m.CountersValuesPassed, value)
	return 0, nil
}

// assertions

func (m *MetricsUpdaterMock) AssertCalled(t *testing.T, times int) {
	t.Helper()
	if m.Calls != times {
		t.Errorf("want store to be called exactly %d times, got %d", times, m.Calls)
	}
}

func (m *MetricsUpdaterMock) AssertNames(t *testing.T, names []string) {
	t.Helper()
	if !slices.Equal(m.NamesPassed, names) {
		t.Errorf("want store to be called with names %v, got %v", names, m.NamesPassed)
	}
}

func (m *MetricsUpdaterMock) AssertGaugesValues(t *testing.T, values []float64) {
	t.Helper()
	if !slices.Equal(m.GaugesValuesPassed, values) {
		t.Errorf("want store to be called with gauges values %v, got %v", values, m.GaugesValuesPassed)
	}
}

func (m *MetricsUpdaterMock) AssertCountersValues(t *testing.T, values []int64) {
	t.Helper()
	if !slices.Equal(m.CountersValuesPassed, values) {
		t.Errorf("want store to be called with counters values %v, got %v", values, m.CountersValuesPassed)
	}
}
