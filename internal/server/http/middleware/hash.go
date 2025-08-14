package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// HeaderHashSHA256 — имя заголовка с HMAC-SHA256 подписью тела.
const HeaderHashSHA256 = "HashSHA256"

type hashWriter struct {
	http.ResponseWriter
	hashKey string
}

func (h *hashWriter) Write(p []byte) (int, error) {
	n, err := h.ResponseWriter.Write(p)
	if err != nil {
		return n, err
	}

	if n > 0 {
		hash := hmac.New(sha256.New, []byte(h.hashKey))
		_, err = hash.Write(p)
		if err != nil {
			return n, err
		}
		h.Header().Set(HeaderHashSHA256, base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	}
	return n, nil
}

// Hash проверяет подпись тела входящего запроса (если заголовок присутствует)
// и добавляет подпись к исходящему ответу. Для вычисления используется
// HMAC-SHA256 с переданным ключом.
func Hash(hashKey string, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// request
			if r.Header.Get(HeaderHashSHA256) != "" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Error(err)
					http.Error(w, "failed to read request body", http.StatusInternalServerError)
					return
				}

				hash := hmac.New(sha256.New, []byte(hashKey))
				_, err = hash.Write(body)
				if err != nil {
					logger.Error(err)
					http.Error(w, "failed to hash request body", http.StatusInternalServerError)
					return
				}

				if r.Header.Get(HeaderHashSHA256) != base64.StdEncoding.EncodeToString(hash.Sum(nil)) {
					msg := fmt.Sprintf("invalid %s header value", HeaderHashSHA256)
					logger.Error(msg)
					http.Error(w, msg, http.StatusBadRequest)
					return
				}

				r.Body = io.NopCloser(bytes.NewReader(body))
			}

			// response writer
			writer := &hashWriter{w, hashKey}

			h.ServeHTTP(writer, r)
		})
	}
}
