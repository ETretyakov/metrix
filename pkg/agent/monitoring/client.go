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

	"metrix/pkg/logger"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client      *resty.Client
	baseURL     string
	signKey     string
	useBatching bool
}

func NewClient(
	ctx context.Context,
	baseURL string,
	signKey string,
	useBatching bool,
	retryCount int,
	retryWaitTime time.Duration,
	retryMaxWaitTime time.Duration,
) *Client {
	c := &Client{
		client:      resty.New(),
		baseURL:     baseURL,
		signKey:     signKey,
		useBatching: useBatching,
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
		return false, fmt.Errorf("failed to request for batching support: %w", err)
	} else if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (c Client) sendMetricBatch(
	ctx context.Context,
	metrics []*Metric,
) error {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	payload, err := json.Marshal(&metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = writer.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	req := c.client.R().
		SetContext(ctx).
		SetBody(&buf)

	if c.signKey != "" {
		h := hmac.New(sha256.New, []byte(c.signKey))
		h.Write(buf.Bytes())
		signature := h.Sum(nil)
		req = req.SetHeader("HashSHA256", hex.EncodeToString(signature))
	}

	resp, err := req.Post(c.baseURL + "/updates/")
	if err != nil {
		return fmt.Errorf("failed to send metric: %w", err)
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

	return nil
}

func (c Client) sendMetric(
	ctx context.Context,
	metrics []*Metric,
) error {
	for _, metric := range metrics {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)

		payload, err := json.Marshal(&metric)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		_, err = writer.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to compress payload: %w", err)
		}

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		req := c.client.R().
			SetContext(ctx).
			SetBody(&buf)

		if c.signKey != "" {
			h := hmac.New(sha256.New, []byte(c.signKey))
			h.Write(buf.Bytes())
			signature := h.Sum(nil)
			req = req.SetHeader("HashSHA256", hex.EncodeToString(signature))
		}

		resp, err := req.Post(c.baseURL + "/update/")
		if err != nil {
			return fmt.Errorf("failed to send metric: %w", err)
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

func (c Client) SendMetrics(ctx context.Context, metrics []*Metric) error {
	if c.useBatching {
		err := c.sendMetricBatch(ctx, metrics)
		if err != nil {
			return fmt.Errorf("failed to send batch metrics: %w", err)
		}
	} else {
		err := c.sendMetric(ctx, metrics)
		if err != nil {
			return fmt.Errorf("failed to send metrics: %w", err)
		}
	}

	return nil
}
