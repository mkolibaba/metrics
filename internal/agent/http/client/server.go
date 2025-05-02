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

func (s *ServerClient) UpdateCounter(name string, value int64) error {
	return s.sendMetric(&model.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &value,
	})
}

func (s *ServerClient) UpdateGauge(name string, value float64) error {
	return s.sendMetric(&model.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	})
}

func (s *ServerClient) sendMetric(body *model.Metrics) error {
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
		Post("/update/")
	return err
}
