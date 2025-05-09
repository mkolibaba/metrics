package mocks

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

func (m *MetricsStorageMock) UpdateGauge(name string, value float64) {
	m.Calls++
	m.NamesPassed = append(m.NamesPassed, name)
	m.GaugesValuesPassed = append(m.GaugesValuesPassed, value)
}

func (m *MetricsStorageMock) UpdateCounter(name string, value int64) {
	m.Calls++
	m.NamesPassed = append(m.NamesPassed, name)
	m.CountersValuesPassed = append(m.CountersValuesPassed, value)
}
