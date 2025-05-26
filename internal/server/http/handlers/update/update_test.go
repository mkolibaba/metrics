package update

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateHandlerShouldReturnCorrectStatus(t *testing.T) {
	cases := []struct {
		url        string
		wantStatus int
	}{
		{
			"/update/counter/my/12",
			200,
		}, {
			"/update/counter/my/-12",
			200,
		}, {
			"/update/gauge/my/1.2",
			200,
		}, {
			"/update/gauge/my/12",
			200,
		}, {
			"/update/counter/123",
			404,
		}, {
			"/update/gauge/123",
			404,
		}, {
			"/update/counter/abc/aaa",
			400,
		}, {
			"/update/counter/abc/1.2",
			400,
		}, {
			"/update/counter/abc/9999999999999999999999999999999999999999999999999999999999",
			400,
		}, {
			"/update/gauge/abc/blabla",
			400,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("POST_%s_should_return_response_with_status_%d", c.url, c.wantStatus), func(t *testing.T) {
			store := &MetricsUpdaterMock{}
			response := sendUpdateRequest(t, store, c.url)

			result := response.Result()
			defer result.Body.Close()
			testutils.AssertResponseStatusCode(t, c.wantStatus, result.StatusCode)
		})
	}
}

func TestUpdateHandlerCallsStoreCorrectly(t *testing.T) {
	type want struct {
		calls                int
		namesPassed          []string
		gaugesValuesPassed   []float64
		countersValuesPassed []int64
	}

	assertState := func(t *testing.T, store *MetricsUpdaterMock, want want) {
		store.AssertCalled(t, want.calls)
		store.AssertNames(t, want.namesPassed)
		store.AssertGaugesValues(t, want.gaugesValuesPassed)
		store.AssertCountersValues(t, want.countersValuesPassed)
	}

	t.Run("Should_call_store_exactly_1_time", func(t *testing.T) {
		store := &MetricsUpdaterMock{}

		sendUpdateRequest(t, store, "/update/counter/my/12")

		assertState(t, store, want{1, []string{"my"}, []float64{}, []int64{12}})
	})
	t.Run("Should_call_store_exactly_2_times", func(t *testing.T) {
		store := &MetricsUpdaterMock{}

		sendUpdateRequest(t, store, "/update/counter/my/12")
		sendUpdateRequest(t, store, "/update/counter/my/12")

		assertState(t, store, want{2, []string{"my", "my"}, []float64{}, []int64{12, 12}})
	})
	t.Run("Should_correctly_process_all_requests", func(t *testing.T) {
		store := &MetricsUpdaterMock{}

		sendUpdateRequest(t, store, "/update/counter/a/1")
		sendUpdateRequest(t, store, "/update/counter/b/5")
		sendUpdateRequest(t, store, "/update/gauge/d/3")
		sendUpdateRequest(t, store, "/update/gauge/e/4")
		sendUpdateRequest(t, store, "/update/counter/c/2")

		assertState(t, store, want{5, []string{"a", "b", "d", "e", "c"}, []float64{3, 4}, []int64{1, 5, 2}})
	})
	t.Run("Should_not_call_store_when_request_is_invalid", func(t *testing.T) {
		store := &MetricsUpdaterMock{}

		sendUpdateRequest(t, store, "/update/counter/abc/1.2")
		sendUpdateRequest(t, store, "/update/gauge/abc/blabla")

		assertState(t, store, want{0, []string{}, []float64{}, []int64{}})
	})
}

func TestSendMetricJSON(t *testing.T) {
	type want struct {
		status   int
		calls    int
		names    []string
		counters []int64
		gauges   []float64
	}

	cases := []struct {
		name string
		body string
		want want
	}{
		{
			name: "should_update_counter",
			body: "{\"id\": \"counter1\",\"type\": \"counter\",\"delta\": 12}",
			want: want{
				status:   200,
				calls:    1,
				names:    []string{"counter1"},
				counters: []int64{12},
			},
		},
		{
			name: "should_update_gauge",
			body: "{\"id\": \"gauge1\",\"type\": \"gauge\",\"value\": 12.34}",
			want: want{
				status: 200,
				calls:  1,
				names:  []string{"gauge1"},
				gauges: []float64{12.34},
			},
		},
		{
			name: "invalid_update_counter",
			body: "{\"id\": \"counter1\",\"type\": \"counter\",\"delta\": 12.34}",
			want: want{
				status: 400,
			},
		},
		{
			name: "invalid_update_gauge",
			body: "{\"id\": \"gauge1\",\"type\": \"gauge\",\"delta\": 12.34}",
			want: want{
				status: 400,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store := &MetricsUpdaterMock{}

			response := sendUpdateRequestJSON(t, store, c.body)

			result := response.Result()
			defer result.Body.Close()
			testutils.AssertResponseStatusCode(t, c.want.status, result.StatusCode)
			store.AssertCalled(t, c.want.calls)
			store.AssertNames(t, c.want.names)
			store.AssertCountersValues(t, c.want.counters)
			store.AssertGaugesValues(t, c.want.gauges)
		})
	}
}

func TestSendMetricResponseJSON(t *testing.T) {
	t.Run("should_return_new_counter", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		body := "{\"id\": \"counter1\",\"type\": \"counter\",\"delta\": 12}"

		response := sendUpdateRequestJSON(t, store, body)
		want := testutils.CreateCounterResponseBodyJSON("counter1", 12)
		testutils.AssertResponseBodyJSON(t, want, response.Body)
	})
	t.Run("should_return_updated_counter", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter(t.Context(), "counter1", 12)
		body := "{\"id\": \"counter1\",\"type\": \"counter\",\"delta\": 12}"

		response := sendUpdateRequestJSON(t, store, body)
		want := testutils.CreateCounterResponseBodyJSON("counter1", 24)
		testutils.AssertResponseBodyJSON(t, want, response.Body)
	})
	t.Run("should_return_new_gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		body := "{\"id\": \"gauge1\",\"type\": \"gauge\",\"value\": 12.34}"

		response := sendUpdateRequestJSON(t, store, body)
		want := testutils.CreateGaugeResponseBodyJSON("gauge1", 12.34)
		testutils.AssertResponseBodyJSON(t, want, response.Body)
	})
	t.Run("should_return_updated_gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge(t.Context(), "gauge1", 12.34)
		body := "{\"id\": \"gauge1\",\"type\": \"gauge\",\"value\": 12.34}"

		response := sendUpdateRequestJSON(t, store, body)
		want := testutils.CreateGaugeResponseBodyJSON("gauge1", 12.34)
		testutils.AssertResponseBodyJSON(t, want, response.Body)
	})
}

func sendUpdateRequestJSON(t *testing.T, updater MetricsUpdater, body string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	server := testutils.NewTestServer("POST /update/", NewJSON(updater, zap.S()))
	return server.Execute(request)
}

func sendUpdateRequest(t *testing.T, updater MetricsUpdater, url string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, url, nil)
	server := testutils.NewTestServer("POST /update/{type}/{name}/{value}", New(updater))
	return server.Execute(request)
}
