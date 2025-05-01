package testutils

import (
	"encoding/json"
	"maps"
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
