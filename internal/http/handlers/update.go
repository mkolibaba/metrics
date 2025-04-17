package handlers

import (
	"github.com/mkolibaba/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

const (
	MetricGauge   = "gauge"
	MetricCounter = "counter"

	RouteUpdate = "/update/"
)

func NewUpdateHandler(store storage.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pathVariables := strings.Split(strings.TrimPrefix(r.URL.Path, RouteUpdate), "/")
		if len(pathVariables) != 3 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		t, name, val := pathVariables[0], pathVariables[1], pathVariables[2]

		switch t {
		case MetricGauge:
			v, err := strconv.ParseFloat(val, 64)
			if err == nil {
				store.UpdateGauge(name, v)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		case MetricCounter:
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
