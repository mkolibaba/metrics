package updatejson

import (
	"encoding/json"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"net/http"
)

type metricsGetterUpdater interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	// TODO: может из апдейта вовзращать новое значение?
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
}

func New(getterUpdater metricsGetterUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

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
			getterUpdater.UpdateGauge(requestBody.ID, *requestBody.Value)
			val, err := getterUpdater.GetGauge(requestBody.ID)
			if err != nil {
				// TODO: что-то странное
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			responseBody.Value = &val
		case handlers.MetricCounter:
			getterUpdater.UpdateCounter(requestBody.ID, *requestBody.Delta)
			val, err := getterUpdater.GetCounter(requestBody.ID)
			if err != nil {
				// TODO: что-то странное
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			responseBody.Delta = &val
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
