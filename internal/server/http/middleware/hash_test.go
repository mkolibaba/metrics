package middleware

import (
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestBodyHashing(t *testing.T) {
	requestBody := `{"test": "qwerty"}`
	cases := map[string]struct {
		hash       string
		wantStatus int
	}{
		"no_hash": {
			"",
			200,
		},
		"valid_hash": {
			"TRnKzt+OY/Vb3TYw1Xk9Sp6Lj+blkid90P+y6i7M0wI=",
			200,
		},
		"invalid_hash": {
			"blabla",
			400,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	hashKey := "random123"

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
			request.Header.Set("HashSHA256", c.hash)
			recorder := httptest.NewRecorder()

			Hash(hashKey, zaptest.NewLogger(t).Sugar())(handler).ServeHTTP(recorder, request)

			if recorder.Code != c.wantStatus {
				t.Fatalf("want status %d, got %d", c.wantStatus, recorder.Code)
			}
		})
	}
}

func TestResponseBodyHashing(t *testing.T) {
	newHandler := func(body string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, body)
			testutils.AssertNoError(t, err)
		}
	}
	hashKey := "random123"

	cases := map[string]struct {
		body string
		hash string
	}{
		"no_body": {
			"",
			"",
		},
		"has_body": {
			`{"test": "qwerty"}`,
			"TRnKzt+OY/Vb3TYw1Xk9Sp6Lj+blkid90P+y6i7M0wI=",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()

			handler := newHandler(c.body)
			Hash(hashKey, zaptest.NewLogger(t).Sugar())(handler).ServeHTTP(recorder, request)

			hashHeader := recorder.Header().Get("HashSHA256")
			if hashHeader != c.hash {
				t.Fatalf("want hash header %s, got %s", c.hash, hashHeader)
			}
		})
	}
}
