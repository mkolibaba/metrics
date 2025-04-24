package collector

import (
	"math/rand"
	"runtime"
	"time"
)

type MetricsCollector struct {
	gauges       map[string]float64
	counters     map[string]int64
	iterations   int
	pollInterval time.Duration
}

func NewMetricsCollector(pollInterval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		gauges:       make(map[string]float64),
		counters:     make(map[string]int64),
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

func (m *MetricsCollector) GetGauges() map[string]float64 {
	return m.gauges
}

func (m *MetricsCollector) GetCounters() map[string]int64 {
	return m.counters
}

func (m *MetricsCollector) collect() {
	m.iterations++
	stats := getMemStats()
	m.gauges["Alloc"] = float64(stats.Alloc)
	m.gauges["BuckHashSys"] = float64(stats.BuckHashSys)
	m.gauges["Frees"] = float64(stats.Frees)
	m.gauges["GCCPUFraction"] = stats.GCCPUFraction
	m.gauges["GCSys"] = float64(stats.GCSys)
	m.gauges["HeapAlloc"] = float64(stats.HeapAlloc)
	m.gauges["HeapIdle"] = float64(stats.HeapIdle)
	m.gauges["HeapInuse"] = float64(stats.HeapInuse)
	m.gauges["HeapObjects"] = float64(stats.HeapObjects)
	m.gauges["HeapReleased"] = float64(stats.HeapReleased)
	m.gauges["HeapSys"] = float64(stats.HeapSys)
	m.gauges["LastGC"] = float64(stats.LastGC)
	m.gauges["Lookups"] = float64(stats.Lookups)
	m.gauges["MCacheInuse"] = float64(stats.MCacheInuse)
	m.gauges["MCacheSys"] = float64(stats.MCacheSys)
	m.gauges["MSpanInuse"] = float64(stats.MSpanInuse)
	m.gauges["MSpanSys"] = float64(stats.MSpanSys)
	m.gauges["Mallocs"] = float64(stats.Mallocs)
	m.gauges["NextGC"] = float64(stats.NextGC)
	m.gauges["NumForcedGC"] = float64(stats.NumForcedGC)
	m.gauges["NumGC"] = float64(stats.NumGC)
	m.gauges["OtherSys"] = float64(stats.OtherSys)
	m.gauges["PauseTotalNs"] = float64(stats.PauseTotalNs)
	m.gauges["StackInuse"] = float64(stats.StackInuse)
	m.gauges["StackSys"] = float64(stats.StackSys)
	m.gauges["Sys"] = float64(stats.Sys)
	m.gauges["TotalAlloc"] = float64(stats.TotalAlloc)
	m.gauges["RandomValue"] = rand.Float64()
	m.counters["PollCount"] = int64(m.iterations)
}

func getMemStats() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}
