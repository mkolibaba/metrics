package update_json_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/mocks"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMetric(t *testing.T) {
	type want struct {
		status   int
		calls    int
		names    []string
		counters []int64
		gauges   []float64
	}

	cases := []struct {
		name string
		body string
		want want
	}{
		{
			name: "should_update_counter",
			body: "{\"id\": \"counter1\",\"type\": \"counter\",\"delta\": 12}",
			want: want{
				status:   200,
				calls:    1,
				names:    []string{"counter1"},
				counters: []int64{12},
			},
		},
		{
			name: "should_update_gauge",
			body: "{\"id\": \"gauge1\",\"type\": \"gauge\",\"value\": 12.34}",
			want: want{
				status: 200,
				calls:  1,
				names:  []string{"gauge1"},
				gauges: []float64{12.34},
			},
		},
		{
			name: "invalid_update_counter",
			body: "{\"id\": \"counter1\",\"type\": \"counter\",\"delta\": 12.34}",
			want: want{
				status: 400,
			},
		},
		{
			name: "invalid_update_gauge",
			body: "{\"id\": \"gauge1\",\"type\": \"gauge\",\"delta\": 12.34}",
			want: want{
				status: 400,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store := &mocks.MetricsStorageMock{}

			response := sendUpdateRequest(t, store, c.body)

			testutils.AssertResponseStatusCode(t, c.want.status, response)
			store.AssertCalled(t, c.want.calls)
			store.AssertNames(t, c.want.names)
			store.AssertCountersValues(t, c.want.counters)
			store.AssertGaugesValues(t, c.want.gauges)
		})
	}
}

// TODO: implement
//func TestResponse(t *testing.T) {
//	type responseBody struct {
//		ID    string  `json:"id"`
//		MType string  `json:"type"`
//		Delta int64   `json:"delta,omitempty"`
//		Value float64 `json:"value,omitempty"`
//	}
//
//	responseBody := responseBody{}
//	err := json.Unmarshal(response.Body(), &responseBody)
//	if err != nil {
//		t.Errorf("error parsing response body: %v", err)
//	}
//	if !reflect.DeepEqual(responseBody, c.want.responseBody) {
//		t.Errorf("did not get correct response body: got %v, want %v", responseBody, c.want.responseBody)
//	}
//}

func sendUpdateRequest(t *testing.T, store router.MetricsStorage, body any) *resty.Response {
	t.Helper()

	srv := httptest.NewServer(router.New(store))
	defer srv.Close()

	request := resty.New().R()
	request.Method = http.MethodPost
	request.URL = srv.URL + "/update"
	request.Body = body
	request.SetHeader("Content-Type", "application/json")

	response, err := request.Send()
	if err != nil {
		t.Fatalf("error when sending request: %v", err)
	}

	return response
}
