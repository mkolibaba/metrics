package middleware

import (
	"compress/gzip"
	"github.com/mkolibaba/metrics/internal/common/logger"
	"io"
	"net/http"
	"strings"
)

var supportsCompressionContentTypes = []string{"text/html", "application/json"}

type gzipWriter struct {
	http.ResponseWriter
	delegate io.Writer
}

func (g *gzipWriter) Write(p []byte) (n int, err error) {
	return g.delegate.Write(p)
}

func Compressor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// response writer
		writer := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && canCompress(r) {
			gw := gzip.NewWriter(w)
			defer gw.Close()

			writer = &gzipWriter{w, gw}
			writer.Header().Set("Content-Encoding", "gzip")
		}

		// request
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Sugared.Error(err)
				http.Error(w, "failed to decompress request body", http.StatusBadRequest)
				return
			}
			defer gr.Close()

			r.Body = gr
		}

		h.ServeHTTP(writer, r)
	})
}

func canCompress(r *http.Request) bool {
	contentType := r.Header.Get("Accept")
	for _, t := range supportsCompressionContentTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}
	return false
}
