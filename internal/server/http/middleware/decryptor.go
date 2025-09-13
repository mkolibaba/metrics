package middleware

import (
	"bytes"
	"github.com/mkolibaba/metrics/internal/common/rsa"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func Decryptor(decryptor *rsa.Decryptor, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error(err)
				http.Error(w, "failed to read request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			decrypted, err := decryptor.Decrypt(body)
			if err != nil {
				logger.Error(err)
				http.Error(w, "failed to decrypt request body", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(decrypted))

			h.ServeHTTP(w, r)
		})
	}
}
