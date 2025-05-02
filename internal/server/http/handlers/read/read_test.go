package read_test

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

type sendRequestFunc func(store router.MetricsStorage, tp, name string) *resty.Response

func TestReadPlain(t *testing.T) {
	sendRequest := func(store router.MetricsStorage, tp, name string) *resty.Response {
		t.Helper()

		srv := httptest.NewServer(router.New(store))
		defer srv.Close()

		client := resty.New().
			SetBaseURL(srv.URL)

		response, err := client.R().
			SetPathParams(map[string]string{
				"t":    tp,
				"name": name,
			}).
			Execute(http.MethodGet, "/value/{t}/{name}")
		testutils.AssertNoError(t, err)

		return response
	}

	doTestRead(t, sendRequest)
}

func TestReadJSON(t *testing.T) {
	sendRequest := func(store router.MetricsStorage, tp, name string) *resty.Response {
		t.Helper()

		srv := httptest.NewServer(router.New(store))
		defer srv.Close()

		client := resty.New().
			SetBaseURL(srv.URL)

		response, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(createRequestBodyJSON(name, tp)).
			Execute(http.MethodPost, "/value/")
		testutils.AssertNoError(t, err)

		return response
	}

	doTestRead(t, sendRequest)
}

func doTestRead(t *testing.T, sendRequest sendRequestFunc) {
	t.Run("Should_return_counter", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter("counter1", 12)

		response := sendRequest(store, "counter", "counter1")

		// TODO: не самый лучший вариант, нужно поправить
		if response.Request.Header.Get("Content-Type") == "application/json" {
			want := createCounterResponseBodyJSON("counter1", 12)
			testutils.AssertResponseBodyJSON(t, want, response)
		} else {
			testutils.AssertResponseBody(t, "12", response)
		}
	})
	t.Run("Should_return_gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)

		response := sendRequest(store, "gauge", "gauge1")

		// TODO: не самый лучший вариант, нужно поправить
		if response.Request.Header.Get("Content-Type") == "application/json" {
			want := createGaugeResponseBodyJSON("gauge1", 34.56)
			testutils.AssertResponseBodyJSON(t, want, response)
		} else {
			testutils.AssertResponseBody(t, "34.56", response)
		}
	})
	t.Run("Should_handle_unexisted_metric", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendRequest(store, "gauge", "gauge2")

		testutils.AssertResponseStatusCode(t, 404, response)
	})
	t.Run("Should_handle_unexisted_metric_type", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendRequest(store, "lolkek", "gauge1")

		testutils.AssertResponseStatusCode(t, 404, response)
	})
}

func createRequestBodyJSON(id, t string) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", id, t)
}

func createGaugeResponseBodyJSON(id string, val float64) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"gauge\", \"value\": %f}", id, val)
}

func createCounterResponseBodyJSON(id string, val int64) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"counter\", \"delta\": %d}", id, val)
}
