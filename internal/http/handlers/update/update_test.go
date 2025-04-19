package update_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/http/router"
	"github.com/mkolibaba/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SpyMetricsStorage struct {
	calls                int
	namesPassed          []string
	gaugesValuesPassed   []float64
	countersValuesPassed []int64
}

func (s *SpyMetricsStorage) GetGauges() map[string]float64 {
	return nil // TODO: реализовать при необходимости
}

func (s *SpyMetricsStorage) GetCounters() map[string]int64 {
	return nil // TODO: реализовать при необходимости
}

func (s *SpyMetricsStorage) GetGauge(name string) (float64, error) {
	return 0, nil // TODO: реализовать при необходимости
}

func (s *SpyMetricsStorage) GetCounter(name string) (int64, error) {
	return 0, nil // TODO: реализовать при необходимости
}

func (m *SpyMetricsStorage) UpdateGauge(name string, value float64) {
	m.calls++
	m.namesPassed = append(m.namesPassed, name)
	m.gaugesValuesPassed = append(m.gaugesValuesPassed, value)
}

func (m *SpyMetricsStorage) UpdateCounter(name string, value int64) {
	m.calls++
	m.namesPassed = append(m.namesPassed, name)
	m.countersValuesPassed = append(m.countersValuesPassed, value)
}

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
		t.Run(fmt.Sprintf("POST %s should return response with status %d", c.url, c.wantStatus), func(t *testing.T) {
			store := &SpyMetricsStorage{}
			response := sendUpdateRequest(t, store, c.url)

			assert.Equal(t, c.wantStatus, response.StatusCode())
		})
	}
}

func TestUpdateHandlerCallsStoreCorrectly(t *testing.T) {
	t.Run("Should call store exactly 1 time", func(t *testing.T) {
		store := &SpyMetricsStorage{}

		sendUpdateRequest(t, store, "/update/counter/my/12")

		assert.Equal(t, 1, store.calls)
		assert.Equal(t, []string{"my"}, store.namesPassed)
		assert.Equal(t, []int64{12}, store.countersValuesPassed)
		assert.Empty(t, store.gaugesValuesPassed)
	})
	t.Run("Should call store exactly 2 times", func(t *testing.T) {
		store := &SpyMetricsStorage{}

		sendUpdateRequest(t, store, "/update/counter/my/12")
		sendUpdateRequest(t, store, "/update/counter/my/12")

		assert.Equal(t, 2, store.calls)
		assert.Equal(t, []string{"my", "my"}, store.namesPassed)
		assert.Equal(t, []int64{12, 12}, store.countersValuesPassed)
		assert.Empty(t, store.gaugesValuesPassed)
	})
	t.Run("Should correctly process all requests", func(t *testing.T) {
		store := &SpyMetricsStorage{}

		sendUpdateRequest(t, store, "/update/counter/a/1")
		sendUpdateRequest(t, store, "/update/counter/b/5")
		sendUpdateRequest(t, store, "/update/gauge/d/3")
		sendUpdateRequest(t, store, "/update/gauge/e/4")
		sendUpdateRequest(t, store, "/update/counter/c/2")

		assert.Equal(t, 5, store.calls)
		assert.Equal(t, []string{"a", "b", "d", "e", "c"}, store.namesPassed)
		assert.Equal(t, []int64{1, 5, 2}, store.countersValuesPassed)
		assert.Equal(t, []float64{3, 4}, store.gaugesValuesPassed)
	})
	t.Run("Should not call store when request is invalid", func(t *testing.T) {
		store := &SpyMetricsStorage{}

		sendUpdateRequest(t, store, "/update/counter/abc/1.2")
		sendUpdateRequest(t, store, "/update/gauge/abc/blabla")

		assert.Empty(t, store.calls)
		assert.Empty(t, store.namesPassed)
		assert.Empty(t, store.countersValuesPassed)
		assert.Empty(t, store.gaugesValuesPassed)
	})
}

func sendUpdateRequest(t *testing.T, store storage.MetricsStorage, url string) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodPost
	request.URL = srv.URL + url

	response, err := request.Send()
	require.NoError(t, err)

	return response
}
