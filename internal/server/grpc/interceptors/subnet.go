package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	"net/netip"
)

func UnarySubnet(subnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unknown, "peer not found")
		}

		addr, err := netip.ParseAddrPort(p.Addr.String())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		ip := net.ParseIP(addr.Addr().String())
		if !subnet.Contains(ip) {
			return nil, status.Error(codes.PermissionDenied, "ip is forbidden")
		}

		return handler(ctx, req)
	}
}
