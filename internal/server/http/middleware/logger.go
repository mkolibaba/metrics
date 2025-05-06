package middleware

import (
	"github.com/mkolibaba/metrics/internal/common/logger"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writerWrapper := loggingResponseWriter{w, &responseData{status: 200}}

		start := time.Now()
		h.ServeHTTP(&writerWrapper, r)
		duration := time.Since(start)

		logger.Sugared.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", writerWrapper.responseData.status,
			"duration", duration,
			"size", writerWrapper.responseData.size,
		)
	})
}
