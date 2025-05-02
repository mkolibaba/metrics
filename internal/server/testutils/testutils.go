package testutils

import (
	"encoding/json"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"testing"
)

type bodyReader interface {
	Body() []byte
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
}

func AssertResponseStatusCode(t *testing.T, want int, got interface {
	StatusCode() int
}) {
	t.Helper()
	if got.StatusCode() != want {
		t.Errorf("did not get correct response status code: want %d, got %d", want, got.StatusCode())
	}
}

func AssertResponseBody(t *testing.T, want string, got bodyReader) {
	t.Helper()
	body := string(got.Body())
	if body != want {
		t.Errorf("did not get correct response body: want '%s' got '%s'", want, body)
	}
}

func AssertResponseBodyJSON(t *testing.T, want string, got bodyReader) {
	t.Helper()
	if len(got.Body()) == 0 && len(want) == 0 {
		return
	}
	gotMap := make(map[string]any)
	if err := json.Unmarshal(got.Body(), &gotMap); err != nil {
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
