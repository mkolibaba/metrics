package list_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("Should_process_empty_store", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		response := sendRequest(t, store)

		testutils.AssertResponseStatusCode(t, 200, response)
		testutils.AssertResponseBody(t, "", response)
	})
	t.Run("Should_return_list_of_metrics", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter("counter1", 12)
		store.UpdateGauge("gauge1", 34.56)

		response := sendRequest(t, store)

		want := "gauge1: 34.560\ncounter1: 12"

		testutils.AssertResponseBody(t, want, response)

		wantContentType := "text/html"
		gotContentType := response.Header().Get("Content-Type")

		if !strings.Contains(gotContentType, wantContentType) {
			t.Errorf("did not get correct response Content-Type: want '%s' got '%s'", wantContentType, gotContentType)
		}
	})
}

func sendRequest(t *testing.T, store router.MetricsStorage) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodGet
	request.URL = srv.URL + "/"

	response, err := request.Send()
	testutils.AssertNoError(t, err)

	return response
}
