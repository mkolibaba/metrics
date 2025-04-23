package client

import (
	"github.com/go-resty/resty/v2"
	"strconv"
)

type ServerClient struct {
	serverAddress string
}

func New(serverAddress string) *ServerClient {
	return &ServerClient{serverAddress}
}

func (s *ServerClient) UpdateCounter(name string, value int64) error {
	return s.sendMetric("counter", name, strconv.FormatInt(value, 10))
}

func (s *ServerClient) UpdateGauge(name string, value float64) error {
	return s.sendMetric("gauge", name, strconv.FormatFloat(value, 'f', 3, 64))
}

func (s *ServerClient) sendMetric(t, name, val string) error {
	_, err := resty.New().R().
		SetPathParams(map[string]string{
			"t":    t,
			"name": name,
			"val":  val,
		}).
		Post("http://" + s.serverAddress + "/update/{t}/{name}/{val}")
	return err
}
