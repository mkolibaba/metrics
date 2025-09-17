package middleware

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentType(t *testing.T) {
	cases := []struct {
		allowed            []string
		requestContentType string
		wantStatus         int
	}{
		{
			[]string{"application/json"},
			"application/json",
			http.StatusOK,
		},
		{
			[]string{"application/json", "application/xml"},
			"application/xml",
			http.StatusOK,
		},
		{
			[]string{"application/json", "application/xml"},
			"text/plain",
			http.StatusUnsupportedMediaType,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			mw := ContentType(c.allowed...)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Content-Type", c.requestContentType)
			recorder := httptest.NewRecorder()

			mw(testutils.EmptyHTTPHandler).ServeHTTP(recorder, request)

			testutils.AssertResponseStatusCode(t, c.wantStatus, recorder.Code)
		})
	}
}
