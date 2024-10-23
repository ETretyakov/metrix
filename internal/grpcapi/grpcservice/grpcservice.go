package grpcservice

import (
	"context"
	"metrix/internal/closer"
	pb "metrix/internal/grpcapi/proto/v1"
	"metrix/internal/model"
	"metrix/internal/repository"
	"metrix/pkg/logger"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
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

func (gs *GServiceServer) Start(ctx context.Context) {
	gServer := grpc.NewServer()

	go func() {
		logger.Info(ctx, "starting listening http srv at :9090")

		tcpListen, err := net.Listen("tcp", ":9090")
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
