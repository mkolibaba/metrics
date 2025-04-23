package mocks

type CollectorMock struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (c *CollectorMock) StartCollect() {
	// do nothing
}

func (c *CollectorMock) GetGauges() map[string]float64 {
	return c.Gauges
}

func (c *CollectorMock) GetCounters() map[string]int64 {
	return c.Counters
}
