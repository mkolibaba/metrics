// Package client предоставляет HTTP‑клиент для отправки батчей метрик
// на сервер метрик.
package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/mkolibaba/metrics/internal/agent/config"
	"github.com/mkolibaba/metrics/internal/agent/http/client/middleware"
	"github.com/mkolibaba/metrics/internal/common/rsa"
	"net"

	"github.com/go-resty/resty/v2"
	"github.com/mkolibaba/metrics/internal/common/http/model"
	"github.com/mkolibaba/metrics/internal/common/retry"
	"go.uber.org/zap"
)

// ServerClient включает в себя HTTP‑клиент и параметры запросов
// для взаимодействия с сервером метрик.
type ServerClient struct {
	client    *resty.Client
	encryptor rsa.Encryptor
	localIP   net.IP
	logger    *zap.SugaredLogger
}

// New создаёт новый клиент для отправки метрик на указанный адрес сервера.
// serverAddress - адрес сервера, hashKey — ключ хедера HashSHA256,
// encryptor - шифровальщик тела запроса, logger — логгер.
func New(cfg *config.AgentConfig, logger *zap.SugaredLogger) (*ServerClient, error) {
	client := resty.New().
		SetBaseURL("http://" + cfg.ServerAddress).
		SetLogger(logger)

	if cfg.CryptoKey != "" {
		encryptor, err := rsa.NewEncryptor(cfg.CryptoKey)
		if err != nil {
			logger.Fatalf("error creating rsa encryptor: %v", err)
		}
		client.OnBeforeRequest(middleware.Encryptor(encryptor))
	}

	if cfg.Key != "" {
		client.OnBeforeRequest(middleware.Hash(cfg.Key))
	}

	localIP, err := getLocalIP()
	if err != nil {
		return nil, err
	}
	return &ServerClient{
		client:  client,
		localIP: localIP,
		logger:  logger,
	}, nil
}

// UpdateCounters отправляет на сервер counter метрики.
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

// UpdateGauges отправляет на сервер gauge метрик.
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
		SetHeader("X-Real-IP", s.localIP.String()).
		SetBody(requestBody.Bytes())

	return retry.Do(func() error {
		_, err := request.Post("/updates/")
		return err
	})
}

func getLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP, nil
}
