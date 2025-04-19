package read_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/http/router"
	"github.com/mkolibaba/metrics/internal/storage"
	"github.com/mkolibaba/metrics/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRead(t *testing.T) {
	t.Run("Should return counter", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, "/value/counter/counter1")

		want := "12"
		got := string(response.Body())

		assert.Equal(t, want, got)
	})
	t.Run("Should return gauge", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)

		response := sendReadRequest(t, store, "/value/gauge/gauge1")

		want := "34.56"
		got := string(response.Body())

		assert.Equal(t, want, got)
	})
	t.Run("Should handle unexisted metric", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, "/value/gauge/gauge2")

		assert.Equal(t, 404, response.StatusCode())
	})
	t.Run("Should handle unexisted metric type", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateGauge("gauge1", 34.56)
		store.UpdateCounter("counter1", 12)

		response := sendReadRequest(t, store, "/value/lolkek/gauge1")

		assert.Equal(t, 404, response.StatusCode())
	})
}

func sendReadRequest(t *testing.T, store storage.MetricsStorage, url string) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodGet
	request.URL = srv.URL + url

	response, err := request.Send()
	require.NoError(t, err)

	return response
}
