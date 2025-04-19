package sender

import (
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type SpyServerAPI struct {
	counterCalls int
	gaugeCalls   int
}

func (s *SpyServerAPI) UpdateCounter(name string, value int64) error {
	s.counterCalls++
	return nil
}

func (s *SpyServerAPI) UpdateGauge(name string, value float64) error {
	s.gaugeCalls++
	return nil
}

func TestSend(t *testing.T) {
	c := collector.NewMetricsCollector(1 * time.Second)
	serverAPI := &SpyServerAPI{}
	sender := NewMetricsSender(c, serverAPI, 1*time.Second)

	c.Gauges["gauge1"] = 1.2
	c.Counters["counter1"] = 2
	c.Counters["counter2"] = 3

	sender.send()

	assert.Equal(t, 2, serverAPI.counterCalls)
	assert.Equal(t, 1, serverAPI.gaugeCalls)
}
