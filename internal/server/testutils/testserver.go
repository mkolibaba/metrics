package testutils

import (
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	h http.Handler
}

func NewTestServer(h http.Handler) *TestServer {
	return &TestServer{
		h: h,
	}
}

func (s *TestServer) Execute(r *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	s.h.ServeHTTP(recorder, r)
	return recorder
}
