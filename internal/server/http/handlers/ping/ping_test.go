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

func TestHandlePing(t *testing.T) {
	cases := map[string]struct {
		pingFunc   func(ctx context.Context) error
		wantStatus int
	}{
		"success": {
			pingFunc: func(ctx context.Context) error {
				return nil
			},
			wantStatus: http.StatusOK,
		},
		"failure": {
			pingFunc: func(ctx context.Context) error {
				return sql.ErrConnDone
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			mock := &mockDB{
				pingFunc: c.pingFunc,
			}

			server := testutils.NewTestServer("/ping", New(mock))

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			response := server.Execute(request)

			testutils.AssertResponseStatusCode(t, c.wantStatus, response.Code)
		})
	}
}
