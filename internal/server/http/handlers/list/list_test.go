package list

import (
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("Should_process_empty_store", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		response := sendRequest(t, store)

		result := response.Result()
		defer result.Body.Close()
		testutils.AssertResponseStatusCode(t, 200, result.StatusCode)
		testutils.AssertResponseBody(t, "<!DOCTYPE html><html><body></body></html>", result.Body)
	})
	t.Run("Should_return_list_of_metrics", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter(t.Context(), "counter1", 12)
		store.UpdateGauge(t.Context(), "gauge1", 34.56)

		response := sendRequest(t, store)

		want := "<!DOCTYPE html><html><body>gauge1: 34.560<br>counter1: 12</body></html>"

		testutils.AssertResponseBody(t, want, response.Body)

		wantContentType := "text/html"
		gotContentType := response.Header().Get("Content-Type")

		if !strings.Contains(gotContentType, wantContentType) {
			t.Errorf("did not get correct response Content-Type: want '%s' got '%s'", wantContentType, gotContentType)
		}
	})
}

func sendRequest(t *testing.T, getter AllMetricsGetter) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	server := testutils.NewTestServer("GET /", New(getter, zaptest.NewLogger(t).Sugar()))
	return server.Execute(request)
}
