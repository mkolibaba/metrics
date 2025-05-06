package inmemory

import "github.com/mkolibaba/metrics/internal/server/storage"

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) GetGauges() map[string]float64 {
	return m.gauges
}

func (m *MemStorage) GetCounters() map[string]int64 {
	return m.counters
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	v, ok := m.gauges[name]
	if !ok {
		return 0, storage.ErrMetricNotFound
	}
	return v, nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	v, ok := m.counters[name]
	if !ok {
		return 0, storage.ErrMetricNotFound
	}
	return v, nil
}

func (m *MemStorage) UpdateGauge(name string, value float64) float64 {
	m.gauges[name] = value
	return m.gauges[name]
}

func (m *MemStorage) UpdateCounter(name string, value int64) int64 {
	m.counters[name] += value
	return m.counters[name]
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
