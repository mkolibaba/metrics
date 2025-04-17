package client

import (
	"fmt"
	"net/http"
	"strconv"
)

type ServerApi interface {
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
}

type ServerClient struct{}

func (s *ServerClient) UpdateCounter(name string, value int64) error {
	return sendMetric("counter", name, strconv.FormatInt(value, 10))
}

func (s *ServerClient) UpdateGauge(name string, value float64) error {
	return sendMetric("gauge", name, strconv.FormatFloat(value, 'f', 4, 64))
}

func sendMetric(t, name, val string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", t, name, val)
	resp, err := http.Post(url, "text/plain", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
