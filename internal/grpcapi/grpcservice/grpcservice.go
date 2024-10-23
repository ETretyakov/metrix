package grpcservice

import (
	"context"
	"metrix/internal/closer"
	pb "metrix/internal/grpcapi/proto/v1"
	"metrix/internal/model"
	"metrix/internal/repository"
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

func (gs *GServiceServer) Start() error {
	tcpListen, err := net.Listen("tcp", ":9090")
	if err != nil {
		return errors.Wrap(err, "failed to start listen tcp")
	}

	gServer := grpc.NewServer()
	pb.RegisterMetricServiceServer(gServer, gs)
	if err := gServer.Serve(tcpListen); err != nil {
		return errors.Wrap(err, "failed to serve grpc server")
	}

	closer.Add(
		func() error {
			gServer.Stop()
			return nil
		},
	)

	return nil
}

func (gs *GServiceServer) MetricsRequest(
	ctx context.Context,
	in *pb.MetricsRequest,
) (*pb.MetricsResponse, error) {
	metrics := []model.Metric{}

	for _, m := range in.Items {
		switch m.Mtype {
		case pb.Metric_COUNTER:
			delta := int64(m.Value)
			metrics = append(metrics, model.Metric{
				ID:    m.Id,
				MType: model.CounterType,
				Delta: &delta,
			})
		case pb.Metric_GAUGE:
			value := float64(m.Value)
			metrics = append(metrics, model.Metric{
				ID:    m.Id,
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
