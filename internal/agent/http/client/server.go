package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"go.uber.org/zap"
	"time"
)

const retryAttempts = 3

var retryIntervalsSeconds = []int{1, 3, 5}

type ServerClient struct {
	client *resty.Client
	logger *zap.SugaredLogger
}

func New(serverAddress string, logger *zap.SugaredLogger) *ServerClient {
	client := resty.New().
		SetBaseURL("http://" + serverAddress).
		SetLogger(logger)
	return &ServerClient{
		client: client,
		logger: logger,
	}
}

func (s *ServerClient) UpdateCounters(counters map[string]int64) error {
	metrics := make([]model.Metrics, 0, len(counters))
	for name, delta := range counters {
		metrics = append(metrics, model.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		})
	}
	return s.sendMetric(metrics)
}

func (s *ServerClient) UpdateGauges(gauges map[string]float64) error {
	metrics := make([]model.Metrics, 0, len(gauges))
	for name, value := range gauges {
		metrics = append(metrics, model.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}
	return s.sendMetric(metrics)
}

func (s *ServerClient) sendMetric(body []model.Metrics) error {
	var compressedBody bytes.Buffer
	gw := gzip.NewWriter(&compressedBody)
	if err := json.NewEncoder(gw).Encode(body); err != nil {
		return err
	}
	gw.Close()

	request := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody.Bytes())

	var err error
	for i := 0; i <= retryAttempts; i++ {
		_, err = request.Post("/updates/")
		if err == nil || i == retryAttempts {
			break
		}

		interval := time.Duration(retryIntervalsSeconds[i]) * time.Second
		s.logger.Warnf("send metrics error: %s. retrying in %s", err, interval)
		time.Sleep(interval)
	}

	return err
}
