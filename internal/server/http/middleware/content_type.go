package middleware

import (
	"net/http"
	"strings"
)

// ContentType проверяет заголовок Content-Type входящих запросов и
// допускает только перечисленные значения. В противном случае
// возвращает статус 415 Unsupported Media Type.
func ContentType(allowedContentTypes ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedContentTypes))
	for _, t := range allowedContentTypes {
		allowed[strings.ToLower(t)] = struct{}{}
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := strings.ToLower(r.Header.Get("Content-Type"))

			if _, ok := allowed[contentType]; !ok {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
