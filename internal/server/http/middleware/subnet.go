package middleware

import (
	"net"
	"net/http"
)

func Subnet(n *net.IPNet) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := net.ParseIP(r.Header.Get("X-Real-IP"))
			if ip == nil || !n.Contains(ip) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
