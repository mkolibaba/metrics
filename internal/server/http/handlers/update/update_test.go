package update_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"github.com/mkolibaba/metrics/internal/server/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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
		t.Run(fmt.Sprintf("POST %s should return response with status %d", c.url, c.wantStatus), func(t *testing.T) {
			store := &mocks.MetricsStorageMock{}
			response := sendUpdateRequest(t, store, c.url)

			assert.Equal(t, c.wantStatus, response.StatusCode())
		})
	}
}

func TestUpdateHandlerCallsStoreCorrectly(t *testing.T) {
	t.Run("Should call store exactly 1 time", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/my/12")

		assert.Equal(t, 1, store.Calls)
		assert.Equal(t, []string{"my"}, store.NamesPassed)
		assert.Equal(t, []int64{12}, store.CountersValuesPassed)
		assert.Empty(t, store.GaugesValuesPassed)
	})
	t.Run("Should call store exactly 2 times", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/my/12")
		sendUpdateRequest(t, store, "/update/counter/my/12")

		assert.Equal(t, 2, store.Calls)
		assert.Equal(t, []string{"my", "my"}, store.NamesPassed)
		assert.Equal(t, []int64{12, 12}, store.CountersValuesPassed)
		assert.Empty(t, store.GaugesValuesPassed)
	})
	t.Run("Should correctly process all requests", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/a/1")
		sendUpdateRequest(t, store, "/update/counter/b/5")
		sendUpdateRequest(t, store, "/update/gauge/d/3")
		sendUpdateRequest(t, store, "/update/gauge/e/4")
		sendUpdateRequest(t, store, "/update/counter/c/2")

		assert.Equal(t, 5, store.Calls)
		assert.Equal(t, []string{"a", "b", "d", "e", "c"}, store.NamesPassed)
		assert.Equal(t, []int64{1, 5, 2}, store.CountersValuesPassed)
		assert.Equal(t, []float64{3, 4}, store.GaugesValuesPassed)
	})
	t.Run("Should not call store when request is invalid", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/abc/1.2")
		sendUpdateRequest(t, store, "/update/gauge/abc/blabla")

		assert.Empty(t, store.Calls)
		assert.Empty(t, store.NamesPassed)
		assert.Empty(t, store.CountersValuesPassed)
		assert.Empty(t, store.GaugesValuesPassed)
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
