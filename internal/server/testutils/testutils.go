package testutils

import (
	"testing"
)

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

func AssertResponseBody(t *testing.T, want string, got interface {
	Body() []byte
}) {
	t.Helper()
	body := string(got.Body())
	if body != want {
		t.Errorf("did not get correct response body: want '%s' got '%s'", want, body)
	}
}
