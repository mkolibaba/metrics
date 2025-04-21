package sender

import (
	"github.com/mkolibaba/metrics/internal/agent/collector"
	"github.com/mkolibaba/metrics/internal/agent/http/client/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	c := collector.NewMetricsCollector(1 * time.Second)
	serverAPI := &mocks.ServerAPIMock{}
	sender := NewMetricsSender(c, serverAPI, 1*time.Second)

	c.Gauges["gauge1"] = 1.2
	c.Counters["counter1"] = 2
	c.Counters["counter2"] = 3

	sender.send()

	assert.Equal(t, 2, serverAPI.CounterCalls)
	assert.Equal(t, 1, serverAPI.GaugeCalls)
}
