package read

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/server/http/handlers"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"go.uber.org/zap"
)

type MetricsGetter interface {
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
}

// New возвращает обработчик, который по пути /value/{type}/{name}
// возвращает значение метрики в текстовом виде. Возвращает 404,
// если метрика не найдена.
func New(getter MetricsGetter, logger *zap.SugaredLogger) http.HandlerFunc {
	writeResponse := func(w http.ResponseWriter, text string) {
		_, err := io.WriteString(w, text)
		if err != nil {
			logger.Errorf("error during processing metrics read request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch t {
		case handlers.MetricCounter:
			counter, err := getter.GetCounter(r.Context(), name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeResponse(w, strconv.FormatInt(counter, 10))
		case handlers.MetricGauge:
			gauge, err := getter.GetGauge(r.Context(), name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeResponse(w, strconv.FormatFloat(gauge, 'f', -1, 64))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// NewJSON возвращает обработчик, ожидающий POST JSON тела с описанием
// метрики и возвращающий её значение в JSON-формате.
func NewJSON(getter MetricsGetter) http.HandlerFunc {
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
		w.Header().Set("Content-Type", "application/json")

		requestBody := &model.Metrics{}
		err := json.NewDecoder(r.Body).Decode(requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t := requestBody.MType
		name := requestBody.ID

		switch t {
		case handlers.MetricCounter:
			counter, err := getter.GetCounter(r.Context(), name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeResponse(w, t, name, &counter, nil)
		case handlers.MetricGauge:
			gauge, err := getter.GetGauge(r.Context(), name)
			if errors.Is(err, storage.ErrMetricNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeResponse(w, t, name, nil, &gauge)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
