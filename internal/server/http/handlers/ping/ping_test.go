package ping

import (
	"context"
	"database/sql"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct {
	pingFunc func(ctx context.Context) error
}

func (m *mockDB) PingContext(ctx context.Context) error {
	if m.pingFunc != nil {
		return m.pingFunc(ctx)
	}
	return nil
}

func TestNewHandler_Success(t *testing.T) {
	mock := &mockDB{
		pingFunc: func(ctx context.Context) error {
			return nil
		},
	}

	server := testutils.NewTestServer("/ping", New(mock))

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := server.Execute(request)

	testutils.AssertResponseStatusCode(t, http.StatusOK, response.Code)
}

func TestNewHandler_Failure(t *testing.T) {
	mock := &mockDB{
		pingFunc: func(ctx context.Context) error {
			return sql.ErrConnDone
		},
	}

	server := testutils.NewTestServer("/ping", New(mock))

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := server.Execute(request)

	testutils.AssertResponseStatusCode(t, http.StatusInternalServerError, response.Code)
}
