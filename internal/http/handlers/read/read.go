package read

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/storage"
	"io"
	"net/http"
	"strconv"
)

func New(store storage.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch t {
		case "counter":
			counter, err := store.GetCounter(name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			io.WriteString(w, strconv.FormatInt(counter, 10))
		case "gauge":
			gauge, err := store.GetGauge(name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			io.WriteString(w, strconv.FormatFloat(gauge, 'f', -1, 64))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
