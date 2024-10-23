package monitoring

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pb "metrix/internal/grpcapi/proto/v1"
	"metrix/pkg/crypto"
	"metrix/pkg/logger"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client - the structure that describes metric client concept.
type Client struct {
	client      *resty.Client
	baseURL     string
	signKey     string
	useBatching bool
	encryption  *crypto.Encryption
}

// NewClient - the builder function for the Client.
func NewClient(
	ctx context.Context,
	baseURL string,
	signKey string,
	useBatching bool,
	retryCount int,
	retryWaitTime time.Duration,
	retryMaxWaitTime time.Duration,
	encryption *crypto.Encryption,
) *Client {
	c := &Client{
		client:      resty.New(),
		baseURL:     baseURL,
		signKey:     signKey,
		useBatching: useBatching,
		encryption:  encryption,
	}

	c.client.
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetRetryMaxWaitTime(retryMaxWaitTime)

	if c.useBatching {
		canBatch, err := c.checkBatching(ctx)
		if err != nil {
			logger.Error(ctx, "failed set batching", err)
			c.useBatching = false
		}

		c.useBatching = canBatch
	}

	return c
}

func (c Client) checkBatching(ctx context.Context) (bool, error) {
	buf := []string{}
	if resp, err := c.client.R().
		SetContext(ctx).
		SetBody(&buf).
		Post(c.baseURL + "/updates/"); err != nil {
		return false, errors.Wrap(err, "failed to request for batching support")
	} else if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (c Client) sendMetricBatch(
	ctx context.Context,
	metrics []*Metric,
) error {
	payload, err := json.Marshal(&metrics)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err = writer.Write(payload); err != nil {
		return errors.Wrap(err, "failed to compress data")
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "failed to close writer")
	}

	var body []byte
	if c.encryption != nil {
		body, err = c.encryption.Encrypt(buf.Bytes())
		if err != nil {
			return errors.Wrap(err, "failed to encrypt body")
		}
	} else {
		body = buf.Bytes()
	}

	req := c.client.R().
		SetContext(ctx).
		SetBody(bytes.NewBuffer(body))

	if c.encryption != nil {
		req = req.SetHeader("x-encrypted", "true")
	}

	if c.signKey != "" {
		h := hmac.New(sha256.New, []byte(c.signKey))
		h.Write(body)
		signature := h.Sum(nil)
		req = req.SetHeader("HashSHA256", hex.EncodeToString(signature))
	}

	resp, err := req.Post(c.baseURL + "/updates/")
	if err != nil {
		return errors.Wrap(err, "failed to send metric")
	}

	if resp.IsError() {
		return errors.Wrapf(
			err,
			"failed  to send metric (with signature): status=%s body=%s",
			resp.Status(),
			resp.Body(),
		)
	}

	logger.Info(
		ctx,
		fmt.Sprintf("sent metric: status=%s body=%s", resp.Status(), resp.Body()),
	)

	return nil
}

func (c Client) sendMetric(
	ctx context.Context,
	metrics []*Metric,
) error {
	for _, metric := range metrics {
		payload, err := json.Marshal(&metric)
		if err != nil {
			return errors.Wrap(err, "failed to marshal payload")
		}

		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)

		_, err = writer.Write(payload)
		if err != nil {
			return errors.Wrap(err, "failed to compress payload")
		}

		err = writer.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close writer")
		}

		var body []byte
		if c.encryption != nil {
			body, err = c.encryption.Encrypt(buf.Bytes())
			if err != nil {
				return errors.Wrap(err, "failed to encrypt body")
			}
		} else {
			body = buf.Bytes()
		}

		req := c.client.R().
			SetContext(ctx).
			SetBody(bytes.NewBuffer(body))

		if c.encryption != nil {
			req = req.SetHeader("x-encrypted", "true")
		}

		if c.signKey != "" {
			h := hmac.New(sha256.New, []byte(c.signKey))
			h.Write(body)
			signature := h.Sum(nil)
			req = req.SetHeader("HashSHA256", hex.EncodeToString(signature))
		}

		resp, err := req.Post(c.baseURL + "/update/")
		if err != nil {
			return errors.Wrap(err, "failed to send metric")
		}

		if resp.IsError() {
			return fmt.Errorf(
				"failed  to send metric (with signature): status=%s body=%s",
				resp.Status(),
				resp.Body(),
			)
		}

		logger.Info(
			ctx,
			fmt.Sprintf("sent metric: status=%s body=%s", resp.Status(), resp.Body()),
		)
	}

	return nil
}

// SendMetrics - the method for sending metrics via metric client.
func (c *Client) SendMetrics(ctx context.Context, metrics []*Metric) error {
	if c.useBatching {
		err := c.sendMetricBatch(ctx, metrics)
		if err != nil {
			return errors.Wrap(err, "failed to send batch metrics")
		}
	} else {
		err := c.sendMetric(ctx, metrics)
		if err != nil {
			return errors.Wrap(err, "failed to send metrics")
		}
	}

	return nil
}

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.MetricServiceClient
}

func NewGRPCClient(serverHost string) *GRPCClient {
	conn, err := grpc.NewClient(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Warn(context.Background(), fmt.Sprintf("failed to build client: %s", err))
	}

	client := pb.NewMetricServiceClient(conn)

	return &GRPCClient{
		conn:   conn,
		client: client,
	}
}

func (gc *GRPCClient) Close() {
	if err := gc.conn.Close(); err != nil {
		logger.Warn(context.Background(), fmt.Sprintf("failed to close client: %s", err))
	}
}

func (gc *GRPCClient) SendMetrics(ctx context.Context, metrics []*Metric) error {
	request := pb.MetricsRequest{}
	for _, m := range metrics {
		switch m.MType {
		case "counter":
			request.Items = append(request.Items, &pb.Metric{
				Id:    m.ID,
				Mtype: pb.Metric_COUNTER,
				Value: float32(*m.Delta),
			})
		case "gauge":
			request.Items = append(request.Items, &pb.Metric{
				Id:    m.ID,
				Mtype: pb.Metric_GAUGE,
				Value: float32(*m.Value),
			})
		}
	}

	resp, err := gc.client.SetMetrics(ctx, &request)
	if err != nil {
		logger.Error(ctx, "failed to send metrics using GRPC", err)
		return errors.Wrap(err, "failed to send metrics using GRPC")
	}

	logger.Info(ctx, fmt.Sprintf("grpc api response: %+v", resp))

	return nil
}
