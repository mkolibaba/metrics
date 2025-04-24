package read_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRead(t *testing.T) {
	t.Run("Should_return_counter", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, "/value/counter/counter1")

		testutils.AssertResponseBody(t, "12", response)
	})
	t.Run("Should_return_gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)

		response := sendReadRequest(t, store, "/value/gauge/gauge1")

		testutils.AssertResponseBody(t, "34.56", response)
	})
	t.Run("Should_handle_unexisted_metric", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, "/value/gauge/gauge2")

		testutils.AssertResponseStatusCode(t, 404, response)
	})
	t.Run("Should_handle_unexisted_metric_type", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, "/value/lolkek/gauge1")

		testutils.AssertResponseStatusCode(t, 404, response)
	})
}

func sendReadRequest(t *testing.T, store router.MetricsStorage, url string) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodGet
	request.URL = srv.URL + url

	response, err := request.Send()
	testutils.AssertNoError(t, err)

	return response
}
