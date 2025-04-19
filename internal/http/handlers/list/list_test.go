package list

import (
	"github.com/mkolibaba/metrics/internal/http/testutils"
	"github.com/mkolibaba/metrics/internal/storage"
	"github.com/mkolibaba/metrics/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("Should process empty store", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		response := sendRequest(store)

		assert.Equal(t, 200, response.StatusCode)
		assert.Empty(t, testutils.ReadAndCloseResponseBody(t, response))
	})
	t.Run("Should return list of metrics", func(t *testing.T) {
		store := inmemory.NewMemStorage()
		store.UpdateCounter("counter1", 12)
		store.UpdateGauge("gauge1", 34.56)

		response := sendRequest(store)

		want := "gauge1: 34.5600\ncounter1: 12"
		got := testutils.ReadAndCloseResponseBody(t, response)

		assert.Equal(t, want, got)
		assert.Contains(t, response.Header.Get("Content-Type"), "text/plain")
	})
}

func sendRequest(store storage.MetricsStorage) *http.Response {
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	recorder := httptest.NewRecorder()

	handler := New(store)
	handler(recorder, request)

	return recorder.Result()
}
