package collector

import (
	"math/rand"
	"runtime"
	"time"
)

type MetricsCollector struct {
	iterations   int
	pollInterval time.Duration
}

func NewMetricsCollector(pollInterval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		pollInterval: pollInterval,
	}
}

func (m *MetricsCollector) StartCollect(chGauges chan<- map[string]float64, chCounters chan<- map[string]int64) {
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			gauges := m.collect()
			chGauges <- gauges
			chCounters <- map[string]int64{"PollCount": int64(m.iterations)}
		}
	}
}

func (m *MetricsCollector) collect() map[string]float64 {
	m.iterations++
	stats := getMemStats()
	gauges := make(map[string]float64)
	gauges["Alloc"] = float64(stats.Alloc)
	gauges["BuckHashSys"] = float64(stats.BuckHashSys)
	gauges["Frees"] = float64(stats.Frees)
	gauges["GCCPUFraction"] = stats.GCCPUFraction
	gauges["GCSys"] = float64(stats.GCSys)
	gauges["HeapAlloc"] = float64(stats.HeapAlloc)
	gauges["HeapIdle"] = float64(stats.HeapIdle)
	gauges["HeapInuse"] = float64(stats.HeapInuse)
	gauges["HeapObjects"] = float64(stats.HeapObjects)
	gauges["HeapReleased"] = float64(stats.HeapReleased)
	gauges["HeapSys"] = float64(stats.HeapSys)
	gauges["LastGC"] = float64(stats.LastGC)
	gauges["Lookups"] = float64(stats.Lookups)
	gauges["MCacheInuse"] = float64(stats.MCacheInuse)
	gauges["MCacheSys"] = float64(stats.MCacheSys)
	gauges["MSpanInuse"] = float64(stats.MSpanInuse)
	gauges["MSpanSys"] = float64(stats.MSpanSys)
	gauges["Mallocs"] = float64(stats.Mallocs)
	gauges["NextGC"] = float64(stats.NextGC)
	gauges["NumForcedGC"] = float64(stats.NumForcedGC)
	gauges["NumGC"] = float64(stats.NumGC)
	gauges["OtherSys"] = float64(stats.OtherSys)
	gauges["PauseTotalNs"] = float64(stats.PauseTotalNs)
	gauges["StackInuse"] = float64(stats.StackInuse)
	gauges["StackSys"] = float64(stats.StackSys)
	gauges["Sys"] = float64(stats.Sys)
	gauges["TotalAlloc"] = float64(stats.TotalAlloc)
	gauges["RandomValue"] = rand.Float64()
	return gauges
}

func getMemStats() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}
