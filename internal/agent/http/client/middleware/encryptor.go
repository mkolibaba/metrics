package middleware

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/common/rsa"
)

func Encryptor(encryptor *rsa.Encryptor) resty.RequestMiddleware {
	return func(client *resty.Client, request *resty.Request) error {
		if body, ok := request.Body.([]byte); ok {
			encrypted, err := encryptor.Encrypt(body)
			if err != nil {
				return err
			}
			request.SetBody(encrypted)
		}
		return nil
	}
}
