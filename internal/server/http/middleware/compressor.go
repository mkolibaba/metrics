package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

var supportedContentTypes = map[string]struct{}{
	"text/html":        {},
	"application/json": {},
}

type gzipWriter struct {
	http.ResponseWriter
	delegate io.Writer
}

func (g *gzipWriter) Write(p []byte) (n int, err error) {
	return g.delegate.Write(p)
}

// Compressor добавляет поддержку gzip-сжатия для ответов и
// распаковку gzip тел входящих запросов. Сжимает только поддерживаемые
// типы содержимого и прозрачно проксирует обработчик.
func Compressor(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
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
					logger.Error(err)
					http.Error(w, "failed to decompress request body", http.StatusInternalServerError)
					return
				}
				defer gr.Close()

				r.Body = gr
			}

			h.ServeHTTP(writer, r)
		})
	}
}

func canCompress(r *http.Request) bool {
	contentType := r.Header.Get("Accept")
	_, ok := supportedContentTypes[contentType]
	return ok
}
