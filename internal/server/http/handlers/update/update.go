package update

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/common/http/model"
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

func NewJSON(updater MetricsUpdater) http.HandlerFunc {
	writeResponse := func(w http.ResponseWriter, t, name string, counter *int64, gauge *float64) {
		responseBody := model.Metrics{
			ID:    name,
			MType: t,
			Delta: counter,
			Value: gauge,
		}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		requestBody := &model.Metrics{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t := requestBody.MType
		name := requestBody.ID

		switch requestBody.MType {
		case handlers.MetricGauge:
			updated := updater.UpdateGauge(requestBody.ID, *requestBody.Value)
			writeResponse(w, t, name, nil, &updated)
		case handlers.MetricCounter:
			updated := updater.UpdateCounter(requestBody.ID, *requestBody.Delta)
			writeResponse(w, t, name, &updated, nil)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
