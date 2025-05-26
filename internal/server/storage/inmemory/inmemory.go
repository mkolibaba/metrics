package inmemory

import (
	"context"
	"github.com/mkolibaba/metrics/internal/server/storage"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) GetGauges(ctx context.Context) (map[string]float64, error) {
	return m.gauges, nil
}

func (m *MemStorage) GetCounters(ctx context.Context) (map[string]int64, error) {
	return m.counters, nil
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	v, ok := m.gauges[name]
	if !ok {
		return 0, storage.ErrMetricNotFound
	}
	return v, nil
}

func (m *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	v, ok := m.counters[name]
	if !ok {
		return 0, storage.ErrMetricNotFound
	}
	return v, nil
}

func (m *MemStorage) UpdateGauge(ctx context.Context, name string, value float64) (float64, error) {
	m.gauges[name] = value
	return m.gauges[name], nil
}

func (m *MemStorage) UpdateCounter(ctx context.Context, name string, value int64) (int64, error) {
	m.counters[name] += value
	return m.counters[name], nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
