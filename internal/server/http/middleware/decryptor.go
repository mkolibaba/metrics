package middleware

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
)

//go:generate moq -stub -out body_decryptor_mock.go . BodyDecryptor
type BodyDecryptor interface {
	Decrypt([]byte) ([]byte, error)
}

func Decryptor(decryptor BodyDecryptor, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
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
