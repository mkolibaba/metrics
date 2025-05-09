package update_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/mocks"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"slices"
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
		t.Run(fmt.Sprintf("POST_%s_should_return_response_with_status_%d", c.url, c.wantStatus), func(t *testing.T) {
			store := &mocks.MetricsStorageMock{}
			response := sendUpdateRequest(t, store, c.url)

			testutils.AssertResponseStatusCode(t, c.wantStatus, response)
		})
	}
}

func TestUpdateHandlerCallsStoreCorrectly(t *testing.T) {
	type want struct {
		calls                int
		namesPassed          []string
		gaugesValuesPassed   []float64
		countersValuesPassed []int64
	}

	assertState := func(t *testing.T, got *mocks.MetricsStorageMock, want want) {
		t.Helper()
		if got.Calls != want.calls {
			t.Errorf("want store to be called exactly %d times, got %d", want.calls, got.Calls)
		}
		if !slices.Equal(got.NamesPassed, want.namesPassed) {
			t.Errorf("want store to be called with names %v, got %v", want.calls, got.Calls)
		}
		if !slices.Equal(got.GaugesValuesPassed, want.gaugesValuesPassed) {
			t.Errorf("want store to be called with gauges values %v, got %v", want.gaugesValuesPassed, got.GaugesValuesPassed)
		}
		if !slices.Equal(got.CountersValuesPassed, want.countersValuesPassed) {
			t.Errorf("want store to be called with counters values %v, got %v", want.countersValuesPassed, got.CountersValuesPassed)
		}
	}

	t.Run("Should_call_store_exactly_1_time", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/my/12")

		assertState(t, store, want{1, []string{"my"}, []float64{}, []int64{12}})
	})
	t.Run("Should_call_store_exactly_2_times", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/my/12")
		sendUpdateRequest(t, store, "/update/counter/my/12")

		assertState(t, store, want{2, []string{"my", "my"}, []float64{}, []int64{12, 12}})
	})
	t.Run("Should_correctly_process_all_requests", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/a/1")
		sendUpdateRequest(t, store, "/update/counter/b/5")
		sendUpdateRequest(t, store, "/update/gauge/d/3")
		sendUpdateRequest(t, store, "/update/gauge/e/4")
		sendUpdateRequest(t, store, "/update/counter/c/2")

		assertState(t, store, want{5, []string{"a", "b", "d", "e", "c"}, []float64{3, 4}, []int64{1, 5, 2}})
	})
	t.Run("Should_not_call_store_when_request_is_invalid", func(t *testing.T) {
		store := &mocks.MetricsStorageMock{}

		sendUpdateRequest(t, store, "/update/counter/abc/1.2")
		sendUpdateRequest(t, store, "/update/gauge/abc/blabla")

		assertState(t, store, want{0, []string{}, []float64{}, []int64{}})
	})
}

func sendUpdateRequest(t *testing.T, store router.MetricsStorage, url string) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodPost
	request.URL = srv.URL + url

	response, err := request.Send()
	if err != nil {
		t.Fatalf("error when sending request: %v", err)
	}

	return response
}
