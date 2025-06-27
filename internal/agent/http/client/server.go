package client

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/common/retry"
	"go.uber.org/zap"
)

type ServerClient struct {
	client  *resty.Client
	hashKey string
	logger  *zap.SugaredLogger
}

func New(serverAddress, hashKey string, logger *zap.SugaredLogger) *ServerClient {
	client := resty.New().
		SetBaseURL("http://" + serverAddress).
		SetLogger(logger)
	return &ServerClient{
		client:  client,
		hashKey: hashKey,
		logger:  logger,
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
	var requestBody bytes.Buffer
	gw := gzip.NewWriter(&requestBody)
	if err := json.NewEncoder(gw).Encode(body); err != nil {
		return err
	}
	gw.Close()

	request := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(requestBody.Bytes())

	if s.hashKey != "" {
		hash := hmac.New(sha256.New, []byte(s.hashKey))
		_, err := hash.Write(requestBody.Bytes())
		if err != nil {
			return err
		}
		request.SetHeader("HashSHA256", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	}

	return retry.Do(func() error {
		_, err := request.Post("/updates/")
		return err
	})
}
