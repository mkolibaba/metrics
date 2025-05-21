package mocks

import (
	"fmt"
)

type ServerAPIMock struct {
	CounterCalls int
	GaugeCalls   int
	ShouldError  bool
}

func (s *ServerAPIMock) UpdateCounter(name string, value int64) error {
	s.CounterCalls++
	if s.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (s *ServerAPIMock) UpdateGauge(name string, value float64) error {
	s.GaugeCalls++
	if s.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}
