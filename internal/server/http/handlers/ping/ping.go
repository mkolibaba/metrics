package ping

import (
	"database/sql"
	"net/http"
)

func New(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.PingContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
