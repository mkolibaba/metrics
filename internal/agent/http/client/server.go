package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/common/http/model"
)

type ServerClient struct {
	serverAddress string
}

func New(serverAddress string) *ServerClient {
	return &ServerClient{serverAddress}
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
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("http://" + s.serverAddress + "/update")
	return err
}
