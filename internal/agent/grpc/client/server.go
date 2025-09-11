package client

import (
	"context"
	metrics "github.com/mkolibaba/metrics/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
)

type ServerClient struct {
	conn *grpc.ClientConn
	c    metrics.ServiceClient
}

func New(serverAddress string) (*ServerClient, error) {
	conn, err := grpc.NewClient(serverAddress)
	if err != nil {
		return nil, err
	}
	c := metrics.NewServiceClient(conn)

	return &ServerClient{
		conn: conn,
		c:    c,
	}, nil
}

func (s *ServerClient) UpdateCounters(counters map[string]int64) error {
	data := make([]*metrics.Metrics, 0, len(counters))
	for name, delta := range counters {
		data = append(data, &metrics.Metrics{
			Id:    name,
			MType: metrics.MType_COUNTER,
			Delta: delta,
		})
	}
	return s.sendMetric(data)
}

func (s *ServerClient) UpdateGauges(gauges map[string]float64) error {
	data := make([]*metrics.Metrics, 0, len(gauges))
	for name, value := range gauges {
		data = append(data, &metrics.Metrics{
			Id:    name,
			MType: metrics.MType_GAUGE,
			Value: value,
		})
	}
	return s.sendMetric(data)
}

func (s *ServerClient) sendMetric(data []*metrics.Metrics) error {
	_, err := s.c.UpdateAll(context.Background(), &metrics.UpdateAllRequest{Data: data})
	return err
}

func (s *ServerClient) Close() {
	s.conn.Close()
}
