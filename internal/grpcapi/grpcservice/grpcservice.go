package grpcservice

import (
	"context"
	"fmt"
	"metrix/internal/closer"
	pb "metrix/internal/grpcapi/proto/v1"
	"metrix/internal/model"
	"metrix/internal/repository"
	"metrix/pkg/logger"
	"net"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type GServiceServer struct {
	pb.UnimplementedMetricServiceServer
	Repository repository.MetricRepository
}

func NewGServiceServer(metricsRepo repository.MetricRepository) *GServiceServer {
	return &GServiceServer{
		Repository: metricsRepo,
	}
}

func subnetInterceptor(
	subnet *net.IPNet,
) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "Failed to get peer from context")
		}

		ipPort := strings.Split(p.Addr.String(), ":")
		ipStr := ipPort[0]
		clientIP := net.ParseIP(ipStr)
		logger.Debug(ctx, fmt.Sprintf("request from ip %s", clientIP))

		if subnet != nil && (clientIP == nil || !subnet.Contains(clientIP)) {
			logger.Warn(ctx, fmt.Sprintf("request from ip %s blocked", clientIP))
			return nil, status.Error(codes.PermissionDenied, "Access denied")
		}

		return handler(ctx, req)
	}
}

func (gs *GServiceServer) Start(ctx context.Context, address string, trustedSubnet *net.IPNet) {
	gServer := grpc.NewServer(grpc.UnaryInterceptor(subnetInterceptor(trustedSubnet)))

	go func() {
		logger.Info(ctx, "starting listening http srv at "+address)

		tcpListen, err := net.Listen("tcp", address)
		if err != nil {
			logger.Fatal(ctx, "failed to start listen tcp", err)
		}

		pb.RegisterMetricServiceServer(gServer, gs)
		if err := gServer.Serve(tcpListen); err != nil {
			logger.Fatal(ctx, "failed to serve grpc server", err)
		}
	}()

	closer.Add(
		func() error {
			gServer.Stop()
			return nil
		},
	)
}

func (gs *GServiceServer) SetMetrics(
	ctx context.Context,
	in *pb.MetricsRequest,
) (*pb.MetricsResponse, error) {
	metrics := []model.Metric{}

	for _, m := range in.GetItems() {
		switch m.GetMtype() {
		case pb.Metric_COUNTER:
			delta := int64(m.GetValue())
			metrics = append(metrics, model.Metric{
				ID:    m.GetId(),
				MType: model.CounterType,
				Delta: &delta,
			})
		case pb.Metric_GAUGE:
			value := float64(m.GetValue())
			metrics = append(metrics, model.Metric{
				ID:    m.GetId(),
				MType: model.CounterType,
				Value: &value,
			})
		default:
			return nil, errors.New("unsupported metric type")
		}
	}

	if _, err := gs.Repository.UpsertMany(ctx, metrics); err != nil {
		return nil, errors.Wrap(err, "failed to upsert many")
	}

	return &pb.MetricsResponse{
		Status:  true,
		Message: "metrics accepted",
	}, nil
}
