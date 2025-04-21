package mocks

type ServerAPIMock struct {
	CounterCalls int
	GaugeCalls   int
}

func (s *ServerAPIMock) UpdateCounter(name string, value int64) error {
	s.CounterCalls++
	return nil
}

func (s *ServerAPIMock) UpdateGauge(name string, value float64) error {
	s.GaugeCalls++
	return nil
}
