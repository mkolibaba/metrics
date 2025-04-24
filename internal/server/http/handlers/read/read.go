package read

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"io"
	"log"
	"net/http"
	"strconv"
)

type MetricsGetter interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

func New(getter MetricsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch t {
		case handlers.MetricCounter:
			counter, err := getter.GetCounter(name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err = io.WriteString(w, strconv.FormatInt(counter, 10))
			if err != nil {
				log.Printf("error during processing metrics read request: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		case handlers.MetricGauge:
			gauge, err := getter.GetGauge(name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err = io.WriteString(w, strconv.FormatFloat(gauge, 'f', -1, 64))
			if err != nil {
				log.Printf("error during processing metrics read request: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
