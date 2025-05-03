package read_test

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type sendRequestFunc func(store router.MetricsStorage, tp, name string) *httptest.ResponseRecorder

func TestReadText(t *testing.T) {
	sendRequest := func(store router.MetricsStorage, tp, name string) *httptest.ResponseRecorder {
		t.Helper()

		url := fmt.Sprintf("/value/%s/%s", tp, name)
		request := httptest.NewRequest(http.MethodGet, url, nil)

		server := testutils.NewTestServer(router.New(store))
		return server.Execute(request)
	}

	doTestRead(t, sendRequest)
}

func TestReadJSON(t *testing.T) {
	sendRequest := func(store router.MetricsStorage, tp, name string) *httptest.ResponseRecorder {
		t.Helper()

		request := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(createRequestBodyJSON(name, tp)))
		request.Header.Set("Content-Type", "application/json")

		server := testutils.NewTestServer(router.New(store))
		return server.Execute(request)
	}

	doTestRead(t, sendRequest)
}

func doTestRead(t *testing.T, sendRequest sendRequestFunc) {
	t.Run("Should_return_counter", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter("counter1", 12)

		response := sendRequest(store, "counter", "counter1")

		// TODO: не самый лучший вариант, нужно поправить
		if response.Header().Get("Content-Type") == "application/json" {
			want := testutils.CreateCounterResponseBodyJSON("counter1", 12)
			testutils.AssertResponseBodyJSON(t, want, response.Body)
		} else {
			testutils.AssertResponseBody(t, "12", response.Body)
		}
	})
	t.Run("Should_return_gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)

		response := sendRequest(store, "gauge", "gauge1")

		// TODO: не самый лучший вариант, нужно поправить
		if response.Header().Get("Content-Type") == "application/json" {
			want := testutils.CreateGaugeResponseBodyJSON("gauge1", 34.56)
			testutils.AssertResponseBodyJSON(t, want, response.Body)
		} else {
			testutils.AssertResponseBody(t, "34.56", response.Body)
		}
	})
	t.Run("Should_handle_unexisted_metric", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendRequest(store, "gauge", "gauge2")

		testutils.AssertResponseStatusCode(t, 404, response.Result().StatusCode)
	})
	t.Run("Should_handle_unexisted_metric_type", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendRequest(store, "lolkek", "gauge1")

		testutils.AssertResponseStatusCode(t, 404, response.Result().StatusCode)
	})
}

func createRequestBodyJSON(id, t string) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"%s\"}", id, t)
}
