package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
)

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func AssertResponseStatusCode(t *testing.T, want int, got int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct response status code: want %d, got %d", want, got)
	}
}

func AssertResponseBody(t *testing.T, want string, got io.Reader) {
	t.Helper()

	if err := iotest.TestReader(got, []byte(want)); err != nil {
		t.Errorf("did not get correct response body: %v", err)
	}
}

func AssertResponseBodyJSON(t *testing.T, want string, got io.Reader) {
	t.Helper()

	bytes, err := io.ReadAll(got)
	AssertNoError(t, err)

	if len(bytes) == 0 && len(want) == 0 {
		return
	}
	gotMap := make(map[string]any)
	if err := json.Unmarshal(bytes, &gotMap); err != nil {
		t.Errorf("error during parse got: %v", err)
	}
	wantMap := make(map[string]any)
	if err := json.NewDecoder(strings.NewReader(want)).Decode(&wantMap); err != nil {
		t.Errorf("error during parse want: %v", err)
	}
	if !maps.Equal(gotMap, wantMap) {
		t.Errorf("did not get correct response body: want '%s' got '%s'", wantMap, gotMap)
	}
}

func CreateGaugeResponseBodyJSON(id string, val float64) string {
	v := strconv.FormatFloat(val, 'f', -1, 64)
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"gauge\", \"value\": %s}", id, v)
}

func CreateCounterResponseBodyJSON(id string, val int64) string {
	return fmt.Sprintf("{\"id\": \"%s\", \"type\": \"counter\", \"delta\": %d}", id, val)
}

// alwaysFailingReader имитирует reader, который всегда возвращает ошибку.
type alwaysFailingReader struct{}

func (er *alwaysFailingReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (er *alwaysFailingReader) Close() error {
	return nil
}

var AlwaysFailingReader = &alwaysFailingReader{}

var EmptyHTTPHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
