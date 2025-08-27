package inmemory

import (
	"context"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"sync"
)

// MemStorage реализует интерфейс storage.Storage и хранит метрики в памяти.
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64

	mu sync.RWMutex
}

// GetGauges возвращает все gauge-метрики из памяти.
func (m *MemStorage) GetGauges(ctx context.Context) (map[string]float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		result[k] = v
	}
	return result, nil
}

// GetCounters возвращает все counter-метрики из памяти.
func (m *MemStorage) GetCounters(ctx context.Context) (map[string]int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		result[k] = v
	}
	return result, nil
}

// GetGauge возвращает значение gauge-метрики по имени.
func (m *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.gauges[name]
	if !ok {
		return 0, storage.ErrMetricNotFound
	}
	return v, nil
}

// GetCounter возвращает значение counter-метрики по имени.
func (m *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.counters[name]
	if !ok {
		return 0, storage.ErrMetricNotFound
	}
	return v, nil
}

// UpdateGauge обновляет значение gauge-метрики.
func (m *MemStorage) UpdateGauge(ctx context.Context, name string, value float64) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
	return m.gauges[name], nil
}

// UpdateCounter обновляет значение counter-метрики.
func (m *MemStorage) UpdateCounter(ctx context.Context, name string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += value
	return m.counters[name], nil
}

// UpdateGauges выполняет батчевое обновление gauge-метрик.
func (m *MemStorage) UpdateGauges(ctx context.Context, values []storage.Gauge) error {
	for _, v := range values {
		m.UpdateGauge(ctx, v.Name, v.Value)
	}
	return nil
}

// UpdateCounters выполняет батчевое обновление counter-метрик.
func (m *MemStorage) UpdateCounters(ctx context.Context, values []storage.Counter) error {
	for _, v := range values {
		m.UpdateCounter(ctx, v.Name, v.Value)
	}
	return nil
}

// NewMemStorage создает новый экземпляр MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
