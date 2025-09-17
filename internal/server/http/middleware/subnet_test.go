package middleware

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubnet(t *testing.T) {
	_, n, err := net.ParseCIDR("192.168.0.0/16")
	testutils.AssertNoError(t, err)

	cases := []struct {
		ip         string
		wantStatus int
	}{
		{
			ip:         "192.168.1.12",
			wantStatus: http.StatusOK,
		},
		{
			ip:         "192.167.1.12",
			wantStatus: http.StatusForbidden,
		},
		{
			ip:         "",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("want_%s_for_", c.ip), func(t *testing.T) {
			mw := Subnet(n)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("X-Real-IP", c.ip)

			recorder := testutils.NewTestServer("GET /", mw(testutils.EmptyHTTPHandler)).Execute(request)

			testutils.AssertResponseStatusCode(t, c.wantStatus, recorder.Code)
		})
	}
}
