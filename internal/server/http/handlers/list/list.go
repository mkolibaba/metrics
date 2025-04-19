package list

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"io"
	"net/http"
	"strings"
)

func New(store storage.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []string
		for k, v := range store.GetGauges() {
			metrics = append(metrics, fmt.Sprintf("%s: %.3f", k, v))
		}
		for k, v := range store.GetCounters() {
			metrics = append(metrics, fmt.Sprintf("%s: %d", k, v))
		}
		_, err := io.WriteString(w, strings.Join(metrics, "\n"))
		if err != nil {
			// никак такую ошибку не обработаем
			panic(err)
		}
	}
}
