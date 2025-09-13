package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/go-resty/resty/v2"
)

func Hash(hashKey string) resty.RequestMiddleware {
	return func(client *resty.Client, request *resty.Request) error {
		if body, ok := request.Body.([]byte); ok {
			hash := hmac.New(sha256.New, []byte(hashKey))
			_, err := hash.Write(body)
			if err != nil {
				return err
			}
			request.SetHeader("HashSHA256", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
		}
		return nil
	}
}
