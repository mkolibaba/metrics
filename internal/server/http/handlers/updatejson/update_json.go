package updatejson

import (
	"encoding/json"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"net/http"
)

type MetricsUpdater interface {
	UpdateGauge(name string, value float64) float64
	UpdateCounter(name string, value int64) int64
}

func New(updater MetricsUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		requestBody := &model.Metrics{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseBody := &model.Metrics{
			ID:    requestBody.ID,
			MType: requestBody.MType,
		}
		switch requestBody.MType {
		case handlers.MetricGauge:
			updated := updater.UpdateGauge(requestBody.ID, *requestBody.Value)
			responseBody.Value = &updated
		case handlers.MetricCounter:
			updated := updater.UpdateCounter(requestBody.ID, *requestBody.Delta)
			responseBody.Delta = &updated
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(responseBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
