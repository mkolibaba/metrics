package collector

import (
	"context"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type MetricsCollector struct {
	iterations   int
	pollInterval time.Duration
	logger       *zap.SugaredLogger
}

func NewMetricsCollector(pollInterval time.Duration, logger *zap.SugaredLogger) *MetricsCollector {
	return &MetricsCollector{
		pollInterval: pollInterval,
		logger:       logger,
	}
}

func (m *MetricsCollector) StartCollect(ctx context.Context) (<-chan map[string]float64, <-chan map[string]int64) {
	chGauges := make(chan map[string]float64, 2)
	chCounters := make(chan map[string]int64, 1)

	go func() {
		ticker := time.NewTicker(m.pollInterval)
		defer ticker.Stop()

	loop:
		for {
			select {
			case <-ticker.C:
				m.iterations++
				chGauges <- m.collectGauges()
				chGauges <- m.collectAdditionalGauges()
				chCounters <- m.collectCounters()
				m.logger.Debug("metrics has been collected")
			case <-ctx.Done():
				m.logger.Debug("stop collecting")
				close(chGauges)
				close(chCounters)
				break loop
			}
		}
	}()

	return chGauges, chCounters
}

func (m *MetricsCollector) collectGauges() map[string]float64 {
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

func (m *MetricsCollector) collectAdditionalGauges() map[string]float64 {
	gauges := make(map[string]float64, 10)

	memory, err := mem.VirtualMemory()
	if err == nil {
		gauges["TotalMemory"] = float64(memory.Total)
		gauges["FreeMemory"] = float64(memory.Free)
	} else {
		m.logger.Errorf("failed to collect additional memory gauges: %s", err)
	}

	usagePerCore, err := cpu.Percent(0, true)
	if err == nil {
		for i, usage := range usagePerCore {
			gauges["CPUutilization"+strconv.Itoa(i+1)] = usage
		}
	} else {
		m.logger.Errorf("failed to collect additional cpu usage gauges: %s", err)
	}

	return gauges
}

func (m *MetricsCollector) collectCounters() map[string]int64 {
	return map[string]int64{
		"PollCount": int64(m.iterations),
	}
}

func getMemStats() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}
