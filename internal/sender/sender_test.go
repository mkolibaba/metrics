package sender

import (
	"github.com/mkolibaba/metrics/internal/collector"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SpyServerApi struct {
	counterCalls int
	gaugeCalls   int
}

func (s *SpyServerApi) UpdateCounter(name string, value int64) error {
	s.counterCalls++
	return nil
}

func (s *SpyServerApi) UpdateGauge(name string, value float64) error {
	s.gaugeCalls++
	return nil
}

func TestSend(t *testing.T) {
	c := collector.NewMetricsCollector()
	serverApi := &SpyServerApi{}
	sender := NewMetricsSender(c, serverApi)

	c.Gauges["gauge1"] = 1.2
	c.Counters["counter1"] = 2
	c.Counters["counter2"] = 3

	sender.send()

	assert.Equal(t, 2, serverApi.counterCalls)
	assert.Equal(t, 1, serverApi.gaugeCalls)
}
