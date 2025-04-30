package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var log *zap.SugaredLogger

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
		writerWrapper := loggingResponseWriter{w, &responseData{}}

		start := time.Now()
		h.ServeHTTP(&writerWrapper, r)
		duration := time.Since(start)

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", writerWrapper.responseData.status,
			"duration", duration,
			"size", writerWrapper.responseData.size,
		)
	})
}

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
}
