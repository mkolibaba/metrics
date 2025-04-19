package collector

import (
	"math/rand"
	"runtime"
	"time"
)

type MetricsCollector struct {
	Gauges       map[string]float64
	Counters     map[string]int64
	iterations   int
	pollInterval time.Duration
}

func NewMetricsCollector(pollInterval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		Gauges:       make(map[string]float64),
		Counters:     make(map[string]int64),
		pollInterval: pollInterval,
	}
}

func (m *MetricsCollector) StartCollect() {
	go func() {
		for {
			m.collect()
			time.Sleep(m.pollInterval)
		}
	}()
}

func (m *MetricsCollector) collect() {
	m.iterations++
	stats := getMemStats()
	m.Gauges["Alloc"] = float64(stats.Alloc)
	m.Gauges["BuckHashSys"] = float64(stats.BuckHashSys)
	m.Gauges["Frees"] = float64(stats.Frees)
	m.Gauges["GCCPUFraction"] = stats.GCCPUFraction
	m.Gauges["GCSys"] = float64(stats.GCSys)
	m.Gauges["HeapAlloc"] = float64(stats.HeapAlloc)
	m.Gauges["HeapIdle"] = float64(stats.HeapIdle)
	m.Gauges["HeapInuse"] = float64(stats.HeapInuse)
	m.Gauges["HeapObjects"] = float64(stats.HeapObjects)
	m.Gauges["HeapReleased"] = float64(stats.HeapReleased)
	m.Gauges["HeapSys"] = float64(stats.HeapSys)
	m.Gauges["LastGC"] = float64(stats.LastGC)
	m.Gauges["Lookups"] = float64(stats.Lookups)
	m.Gauges["MCacheInuse"] = float64(stats.MCacheInuse)
	m.Gauges["MCacheSys"] = float64(stats.MCacheSys)
	m.Gauges["MSpanInuse"] = float64(stats.MSpanInuse)
	m.Gauges["MSpanSys"] = float64(stats.MSpanSys)
	m.Gauges["Mallocs"] = float64(stats.Mallocs)
	m.Gauges["NextGC"] = float64(stats.NextGC)
	m.Gauges["NumForcedGC"] = float64(stats.NumForcedGC)
	m.Gauges["NumGC"] = float64(stats.NumGC)
	m.Gauges["OtherSys"] = float64(stats.OtherSys)
	m.Gauges["PauseTotalNs"] = float64(stats.PauseTotalNs)
	m.Gauges["StackInuse"] = float64(stats.StackInuse)
	m.Gauges["StackSys"] = float64(stats.StackSys)
	m.Gauges["Sys"] = float64(stats.Sys)
	m.Gauges["TotalAlloc"] = float64(stats.TotalAlloc)
	m.Gauges["RandomValue"] = rand.Float64()
	m.Counters["PollCount"] = int64(m.iterations)
}

func getMemStats() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}
