package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/common/http/model"
)

type ServerClient struct {
	client *resty.Client
}

func New(serverAddress string) *ServerClient {
	client := resty.New().
		SetBaseURL("http://" + serverAddress)
	return &ServerClient{
		client: client,
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

	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody.Bytes()).
		Post("/updates/")
	return err
}
