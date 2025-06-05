package update

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"go.uber.org/zap"
)

type MetricsUpdater interface {
	UpdateGauge(ctx context.Context, name string, value float64) (float64, error)
	UpdateCounter(ctx context.Context, name string, value int64) (int64, error)
}

type API struct {
	updater MetricsUpdater
	logger  *zap.SugaredLogger
}

func NewAPI(updater MetricsUpdater, logger *zap.SugaredLogger) *API {
	return &API{
		updater: updater,
		logger:  logger,
	}
}

func (a *API) HandlePlain(w http.ResponseWriter, r *http.Request) {
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
		_, err = a.updater.UpdateGauge(r.Context(), name, v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case handlers.MetricCounter:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = a.updater.UpdateCounter(r.Context(), name, v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (a *API) HandleJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestBody := &model.Metrics{}
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		a.logger.Errorf("can not decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := requestBody.MType
	name := requestBody.ID

	switch t {
	case handlers.MetricGauge:
		updated, err := a.updater.UpdateGauge(r.Context(), name, *requestBody.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		a.writeJSONResponse(w, t, name, nil, &updated)
	case handlers.MetricCounter:
		updated, err := a.updater.UpdateCounter(r.Context(), name, *requestBody.Delta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		a.writeJSONResponse(w, t, name, &updated, nil)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (a *API) HandleJSONBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody []model.Metrics
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		a.logger.Errorf("can not decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, metrics := range requestBody {
		t := metrics.MType
		name := metrics.ID

		switch t {
		case handlers.MetricGauge:
			_, err := a.updater.UpdateGauge(r.Context(), name, *metrics.Value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case handlers.MetricCounter:
			_, err := a.updater.UpdateCounter(r.Context(), name, *metrics.Delta)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func (a *API) writeJSONResponse(w http.ResponseWriter, t, name string, counter *int64, gauge *float64) {
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
