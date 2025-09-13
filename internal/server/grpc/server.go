package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mkolibaba/metrics/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/grpc/interceptors"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"sync"
)

type MetricsStorage interface {
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)

	GetGauges(ctx context.Context) (map[string]float64, error)
	GetCounters(ctx context.Context) (map[string]int64, error)

	UpdateGauge(ctx context.Context, name string, value float64) (float64, error)
	UpdateCounter(ctx context.Context, name string, value int64) (int64, error)

	UpdateGauges(ctx context.Context, values []storage.Gauge) error
	UpdateCounters(ctx context.Context, values []storage.Counter) error
}

type serviceServer struct {
	pb.UnimplementedServiceServer
	store MetricsStorage
}

func (s *serviceServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.Metrics, error) {
	switch in.MType {
	case pb.MType_COUNTER:
		counter, err := s.store.GetCounter(ctx, in.Id)
		if errors.Is(err, storage.ErrMetricNotFound) {
			return nil, status.Error(codes.NotFound, "metric not found")
		}
		return &pb.Metrics{
			Id:    in.Id,
			MType: in.MType,
			Delta: counter,
		}, nil
	case pb.MType_GAUGE:
		gauge, err := s.store.GetGauge(ctx, in.Id)
		if errors.Is(err, storage.ErrMetricNotFound) {
			return nil, status.Error(codes.NotFound, "metric not found")
		}
		return &pb.Metrics{
			Id:    in.Id,
			MType: in.MType,
			Value: gauge,
		}, nil
	default:
		return nil, status.Error(codes.InvalidArgument, "metrics type not supported")
	}
}

func (s *serviceServer) GetAll(ctx context.Context, in *empty.Empty) (*pb.GetAllResponse, error) {
	gauges, err := s.store.GetGauges(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	counters, err := s.store.GetCounters(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Metrics, 0, len(gauges)+len(counters))
	for k, v := range gauges {
		result = append(result, &pb.Metrics{
			Id:    k,
			MType: pb.MType_GAUGE,
			Value: v,
		})
	}
	for k, v := range counters {
		result = append(result, &pb.Metrics{
			Id:    k,
			MType: pb.MType_COUNTER,
			Delta: v,
		})
	}

	return &pb.GetAllResponse{Result: result}, nil
}

func (s *serviceServer) Update(ctx context.Context, in *pb.Metrics) (*pb.Metrics, error) {
	switch in.MType {
	case pb.MType_COUNTER:
		counter, err := s.store.UpdateCounter(ctx, in.Id, in.Delta)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pb.Metrics{
			Id:    in.Id,
			MType: pb.MType_COUNTER,
			Delta: counter,
		}, nil
	case pb.MType_GAUGE:
		gauge, err := s.store.UpdateGauge(ctx, in.Id, in.Value)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pb.Metrics{
			Id:    in.Id,
			MType: pb.MType_GAUGE,
			Value: gauge,
		}, nil
	default:
		return nil, status.Error(codes.InvalidArgument, "metrics type not supported")
	}
}

func (s *serviceServer) UpdateAll(ctx context.Context, in *pb.UpdateAllRequest) (*empty.Empty, error) {
	gauges := make([]storage.Gauge, 0)
	counters := make([]storage.Counter, 0)

	for _, m := range in.Data {
		switch m.MType {
		case pb.MType_COUNTER:
			counters = append(counters, storage.Counter{Name: m.Id, Value: m.Delta})
		case pb.MType_GAUGE:
			gauges = append(gauges, storage.Gauge{Name: m.Id, Value: m.Value})
		}
	}

	if len(gauges) > 0 {
		err := s.store.UpdateGauges(ctx, gauges)
		if err != nil {
			return nil, err
		}
	}

	if len(counters) > 0 {
		err := s.store.UpdateCounters(ctx, counters)
		if err != nil {
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

type Server struct {
	s      *grpc.Server
	logger *zap.SugaredLogger
}

func NewServer(store MetricsStorage, cfg *config.ServerConfig, logger *zap.SugaredLogger) *Server {
	ss := &serviceServer{
		store: store,
	}

	uis := []grpc.UnaryServerInterceptor{interceptors.UnaryLogger(logger)}
	if cfg.TrustedSubnet != nil {
		uis = append(uis, interceptors.UnarySubnet(cfg.TrustedSubnet))
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(uis...),
	)
	pb.RegisterServiceServer(s, ss)
	reflection.Register(s)

	return &Server{
		s:      s,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Infof("running grpc server on %s", addr)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.s.GracefulStop()
	}()

	if err := s.s.Serve(listen); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	s.logger.Infof("server stopped")

	wg.Wait()
	return nil
}
