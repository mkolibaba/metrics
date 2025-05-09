package list

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type AllMetricsGetter interface {
	GetGauges() map[string]float64
	GetCounters() map[string]int64
}

func New(getter AllMetricsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		_, err := io.WriteString(w, strings.Join(metrics, "\n"))
		if err != nil {
			log.Printf("error during processing metrics list request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
