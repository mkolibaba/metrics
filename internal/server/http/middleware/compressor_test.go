package middleware

import (
	"compress/gzip"
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldCompress(t *testing.T) {
	cases := []struct {
		contentType  string
		responseBody string
	}{
		{
			contentType:  "application/json",
			responseBody: "{\"status\": \"ok\"}",
		},
		{
			contentType:  "text/html",
			responseBody: "ok",
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("should_compress_%s", c.contentType), func(t *testing.T) {
			recorder := prepareAndSendRequest(c.contentType, c.responseBody)

			gotContentEncoding := recorder.Header().Get("Content-Encoding")
			if !strings.Contains(gotContentEncoding, "gzip") {
				t.Errorf("error Content-Encoding header: got '%s', want 'gzip'", gotContentEncoding)
			}

			gr, err := gzip.NewReader(recorder.Body)
			testutils.AssertNoError(t, err)

			body, err := io.ReadAll(gr)
			gotBody := string(body)
			testutils.AssertNoError(t, err)

			if gotBody != c.responseBody {
				t.Errorf("error response body: got '%s', want '%s'", gotBody, c.responseBody)
			}
		})
	}
}

func TestShouldNotCompress(t *testing.T) {
	responseBody := "<root>ok</root>"
	contentType := "application/xml"

	recorder := prepareAndSendRequest(contentType, responseBody)

	if recorder.Header().Get("Content-Encoding") != "" {
		t.Errorf("expecting header Content-Encoding to be empty, got %s", recorder.Header().Get("Content-Encoding"))
	}

	gotBody := recorder.Body.String()
	if responseBody != gotBody {
		t.Errorf("error response body: got '%s', want '%s'", gotBody, responseBody)
	}
}

func prepareAndSendRequest(contentType, responseBody string) *httptest.ResponseRecorder {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		io.WriteString(w, responseBody)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept", contentType)
	request.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()

	Compressor(handler).ServeHTTP(recorder, request)

	return recorder
}
