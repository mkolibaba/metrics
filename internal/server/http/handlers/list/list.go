package list

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

const pageTemplate = "<!DOCTYPE html><html><body>%s</body></html>"

type AllMetricsGetter interface {
	GetGauges(ctx context.Context) (map[string]float64, error)
	GetCounters(ctx context.Context) (map[string]int64, error)
}

func New(getter AllMetricsGetter, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		gauges, err := getter.GetGauges(r.Context())
		if err != nil {
			logger.Errorf("error retrieving gauges: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		counters, err := getter.GetCounters(r.Context())
		if err != nil {
			logger.Errorf("error retrieving counters: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		metrics := make([]string, len(gauges)+len(counters))
		var i int
		for k, v := range gauges {
			metrics[i] = fmt.Sprintf("%s: %.3f", k, v)
			i++
		}
		for k, v := range counters {
			metrics[i] = fmt.Sprintf("%s: %d", k, v)
			i++
		}
		_, err = io.WriteString(w, fmt.Sprintf(pageTemplate, strings.Join(metrics, "<br>")))
		if err != nil {
			logger.Errorf("error during processing metrics list request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
