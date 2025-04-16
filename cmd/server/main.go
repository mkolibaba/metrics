package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	MetricGauge   = "gauge"
	MetricCounter = "counter"

	RouteUpdate = "/update/"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// TODO: должно быть место получше
var storage = &MemStorage{
	gauges:   make(map[string]float64),
	counters: make(map[string]int64),
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	trimmed := strings.TrimPrefix(r.URL.Path, RouteUpdate)
	parts := strings.Split(trimmed, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, name, val := parts[0], parts[1], parts[2]

	if !validateMetric(t, val) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t == MetricGauge {
		// игнорируем error, поскольку выше была валидация
		v, _ := strconv.ParseFloat(val, 64)
		storage.gauges[name] = v
	} else if t == MetricCounter {
		// игнорируем error, поскольку выше была валидация
		v, _ := strconv.ParseInt(val, 10, 64)
		storage.counters[name] = v
	}
}

func validateMetric(t, value string) bool {
	switch t {
	case MetricGauge:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case MetricCounter:
		_, err := strconv.ParseInt(value, 10, 64)
		return err == nil
	default:
		return false
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(RouteUpdate, handleUpdate)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
