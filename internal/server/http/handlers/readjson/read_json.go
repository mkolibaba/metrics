package readjson

import (
	"encoding/json"
	"errors"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"net/http"
)

type metricsGetter interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

func New(getter metricsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		requestBody := &model.Metrics{}
		err := json.NewDecoder(r.Body).Decode(requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseBody := model.Metrics{
			ID:    requestBody.ID,
			MType: requestBody.MType,
		}
		switch requestBody.MType {
		case handlers.MetricGauge:
			val, err := getter.GetGauge(requestBody.ID)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			responseBody.Value = &val
		case handlers.MetricCounter:
			val, err := getter.GetCounter(requestBody.ID)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			responseBody.Delta = &val
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err = json.NewEncoder(w).Encode(responseBody); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
