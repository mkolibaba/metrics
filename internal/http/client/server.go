package client

import (
	"github.com/go-resty/resty/v2"
	"strconv"
)

type ServerAPI interface {
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
}

type ServerClient struct{}

func (s *ServerClient) UpdateCounter(name string, value int64) error {
	return sendMetric("counter", name, strconv.FormatInt(value, 10))
}

func (s *ServerClient) UpdateGauge(name string, value float64) error {
	return sendMetric("gauge", name, strconv.FormatFloat(value, 'f', 3, 64))
}

func sendMetric(t, name, val string) error {
	_, err := resty.New().R().
		SetPathParams(map[string]string{
			"t":    t,
			"name": name,
			"val":  val,
		}).
		Post("http://localhost:8080/update/{t}/{name}/{val}")
	return err
}
