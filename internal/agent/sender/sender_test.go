package sender

import (
	cmocks "github.com/mkolibaba/metrics/internal/agent/collector/mocks"
	"github.com/mkolibaba/metrics/internal/agent/http/client/mocks"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	c := &cmocks.CollectorMock{
		Gauges: map[string]float64{
			"gauge1": 1.2,
		},
		Counters: map[string]int64{
			"counter1": 2,
			"counter2": 3,
		},
	}
	serverAPI := &mocks.ServerAPIMock{}
	sender := NewMetricsSender(c, serverAPI, 1*time.Second)

	sender.send()

	want := 2
	got := serverAPI.CounterCalls
	if got != want {
		t.Errorf("want %d counter calls, got %d", want, got)
	}
	want = 1
	got = serverAPI.GaugeCalls
	if got != want {
		t.Errorf("want %d gauge calls, got %d", want, got)
	}
}
