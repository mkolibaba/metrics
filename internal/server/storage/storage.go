package storage

import "errors"

var ErrMetricNotFound = errors.New("metric not found")

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Value int64
}
