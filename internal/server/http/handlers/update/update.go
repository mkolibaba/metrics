package update

import (
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"net/http"
	"strconv"
)

func New(store storage.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		val := chi.URLParam(r, "value")

		switch t {
		case handlers.MetricGauge:
			v, err := strconv.ParseFloat(val, 64)
			if err == nil {
				store.UpdateGauge(name, v)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		case handlers.MetricCounter:
			v, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				store.UpdateCounter(name, v)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
