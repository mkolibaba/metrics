package ping

import (
	"context"
	"net/http"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

// New возвращает обработчик, который проверяет доступность
// внешнего ресурса через Pinger и возвращает 500 при ошибке.
func New(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pinger.PingContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
