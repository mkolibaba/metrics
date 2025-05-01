package read_json_test

import (
	"fmt"
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

		response := sendReadRequest(t, store, createRequestBody("counter1", "counter"))

		want := createCounterResponseBody("counter1", 12)
		testutils.AssertResponseBodyJson(t, want, response)
	})
	t.Run("Should_return_gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)

		response := sendReadRequest(t, store, createRequestBody("gauge1", "gauge"))

		want := createGaugeResponseBody("gauge1", 34.56)
		testutils.AssertResponseBodyJson(t, want, response)
	})
	t.Run("Should_handle_unexisted_metric", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, createRequestBody("gauge2", "gauge"))

		testutils.AssertResponseStatusCode(t, 404, response)
	})
	t.Run("Should_handle_unexisted_metric_type", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, createRequestBody("gauge1", "lolkek"))

		testutils.AssertResponseStatusCode(t, 404, response)
	})
}

func sendReadRequest(t *testing.T, store router.MetricsStorage, body any) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodPost
	request.URL = srv.URL + "/value/"
	request.Body = body
	request.SetHeader("Content-Type", "application/json")

	response, err := request.Send()
	testutils.AssertNoError(t, err)

	return response
}

func createRequestBody(id, t string) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", id, t)
}

func createGaugeResponseBody(id string, val float64) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"gauge\", \"value\": %f}", id, val)
}

func createCounterResponseBody(id string, val int64) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"counter\", \"delta\": %d}", id, val)
}
