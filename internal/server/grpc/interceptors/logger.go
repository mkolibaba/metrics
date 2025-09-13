package interceptors

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func UnaryLogger(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		res, err := handler(ctx, req)
		duration := time.Since(start)

		resStatus, _ := status.FromError(err)
		f := []any{
			"method", info.FullMethod,
			"status", fmt.Sprintf("%d %s", resStatus.Code(), resStatus.Code().String()),
		}
		if err != nil && resStatus.Message() != "" {
			f = append(f, "error", resStatus.Message())
		}
		f = append(f, "duration", duration)

		logger.Infoln(f...)

		return res, err
	}
}
