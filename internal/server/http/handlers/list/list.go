package list

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/common/logger"
	"io"
	"net/http"
	"strings"
)

const pageTemplate = "<!DOCTYPE html><html><body>%s</body></html>"

type AllMetricsGetter interface {
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}

func New(getter AllMetricsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "text/html")
		metrics := make([]string, len(getter.GetGauges())+len(getter.GetCounters()))
		var i int
		for k, v := range getter.GetGauges() {
			metrics[i] = fmt.Sprintf("%s: %.3f", k, v)
			i++
		}
		for k, v := range getter.GetCounters() {
			metrics[i] = fmt.Sprintf("%s: %d", k, v)
			i++
		}
		_, err := io.WriteString(w, fmt.Sprintf(pageTemplate, strings.Join(metrics, "<br>")))
		if err != nil {
			logger.Sugared.Errorf("error during processing metrics list request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
