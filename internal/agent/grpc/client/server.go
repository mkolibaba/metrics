package client

import (
	"context"
	pb "github.com/mkolibaba/metrics/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerClient struct {
	conn *grpc.ClientConn
	c    pb.ServiceClient
}

func New(serverAddress string) (*ServerClient, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewServiceClient(conn)

	return &ServerClient{
		conn: conn,
		c:    c,
	}, nil
}

func (s *ServerClient) UpdateCounters(counters map[string]int64) error {
	data := make([]*pb.Metrics, 0, len(counters))
	for name, delta := range counters {
		m := &pb.Metrics{}
		m.SetId(name)
		m.SetMType(pb.MType_COUNTER)
		m.SetDelta(delta)
		data = append(data, m)
	}
	return s.sendMetric(data)
}

func (s *ServerClient) UpdateGauges(gauges map[string]float64) error {
	data := make([]*pb.Metrics, 0, len(gauges))
	for name, value := range gauges {
		m := &pb.Metrics{}
		m.SetId(name)
		m.SetMType(pb.MType_COUNTER)
		m.SetValue(value)
		data = append(data, m)
	}
	return s.sendMetric(data)
}

func (s *ServerClient) sendMetric(data []*pb.Metrics) error {
	in := &pb.UpdateAllRequest{}
	in.SetData(data)
	_, err := s.c.UpdateAll(context.Background(), in)
	return err
}

func (s *ServerClient) Close() {
	s.conn.Close()
}
