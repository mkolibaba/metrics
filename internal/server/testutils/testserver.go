package testutils

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	router chi.Router
}

func NewTestServer(pattern string, h http.Handler) *TestServer {
	router := chi.NewRouter()
	router.Handle(pattern, h)
	return &TestServer{
		router: router,
	}
}

func (s *TestServer) Execute(r *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	s.router.ServeHTTP(recorder, r)
	return recorder
}
