package storage

import "errors"

var ErrMetricNotFound = errors.New("metric not found")

type MetricsStorage interface {
	GetGauges() map[string]float64
	GetCounters() map[string]int64
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
}
