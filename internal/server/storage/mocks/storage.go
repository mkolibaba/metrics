package mocks

import (
	"slices"
	"testing"
)

type MetricsStorageMock struct {
	Calls                int
	NamesPassed          []string
	GaugesValuesPassed   []float64
	CountersValuesPassed []int64
}

func (m *MetricsStorageMock) GetGauges() map[string]float64 {
	return nil // TODO: реализовать при необходимости
}

func (m *MetricsStorageMock) GetCounters() map[string]int64 {
	return nil // TODO: реализовать при необходимости
}

func (m *MetricsStorageMock) GetGauge(name string) (float64, error) {
	return 0, nil // TODO: реализовать при необходимости
}

func (m *MetricsStorageMock) GetCounter(name string) (int64, error) {
	return 0, nil // TODO: реализовать при необходимости
}

func (m *MetricsStorageMock) UpdateGauge(name string, value float64) float64 {
	m.Calls++
	m.NamesPassed = append(m.NamesPassed, name)
	m.GaugesValuesPassed = append(m.GaugesValuesPassed, value)
	return 0
}

func (m *MetricsStorageMock) UpdateCounter(name string, value int64) int64 {
	m.Calls++
	m.NamesPassed = append(m.NamesPassed, name)
	m.CountersValuesPassed = append(m.CountersValuesPassed, value)
	return 0
}

// assertions

func (m *MetricsStorageMock) AssertCalled(t *testing.T, times int) {
	t.Helper()
	if m.Calls != times {
		t.Errorf("want store to be called exactly %d times, got %d", times, m.Calls)
	}
}

func (m *MetricsStorageMock) AssertNames(t *testing.T, names []string) {
	t.Helper()
	if !slices.Equal(m.NamesPassed, names) {
		t.Errorf("want store to be called with names %v, got %v", names, m.NamesPassed)
	}
}

func (m *MetricsStorageMock) AssertGaugesValues(t *testing.T, values []float64) {
	t.Helper()
	if !slices.Equal(m.GaugesValuesPassed, values) {
		t.Errorf("want store to be called with gauges values %v, got %v", values, m.GaugesValuesPassed)
	}
}

func (m *MetricsStorageMock) AssertCountersValues(t *testing.T, values []int64) {
	t.Helper()
	if !slices.Equal(m.CountersValuesPassed, values) {
		t.Errorf("want store to be called with counters values %v, got %v", values, m.CountersValuesPassed)
	}
}
