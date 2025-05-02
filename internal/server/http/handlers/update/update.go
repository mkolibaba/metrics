package update

import (
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"net/http"
	"strconv"
)

type MetricsUpdater interface {
	UpdateGauge(name string, value float64) float64
	UpdateCounter(name string, value int64) int64
}

func New(updater MetricsUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		val := chi.URLParam(r, "value")

		switch t {
		case handlers.MetricGauge:
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			updater.UpdateGauge(name, v)
		case handlers.MetricCounter:
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			updater.UpdateCounter(name, v)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
