package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"metrix/pkg/logger"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	RetryCount       int           = 3
	RetryWaitTime    time.Duration = time.Second * 2
	RetryMaxWaitTime time.Duration = time.Second * 8
)

type WidgetType string

const (
	CounterType WidgetType = "counter"
	GaugeType   WidgetType = "gauge"
	UnknownType WidgetType = "unknown"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func CheckBatching(
	ctx context.Context,
	baseURL string,
) (bool, error) {
	url := fmt.Sprintf("%s/updates/", baseURL)
	client := resty.New()

	client.
		SetHeader("Accept-Encoding", "gzip").
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime).
		SetRetryMaxWaitTime(RetryMaxWaitTime)

	emptyBuffer := []string{}
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(&emptyBuffer).
		Post(url)
	if err != nil {
		return false, fmt.Errorf("client.ReadAll: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func SendMetricBatch(
	ctx context.Context,
	baseURL string,
	metrics []*Metrics,
) error {
	url := fmt.Sprintf("%s/updates/", baseURL)
	client := resty.New()

	client.
		SetHeader("Accept-Encoding", "gzip").
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime).
		SetRetryMaxWaitTime(RetryMaxWaitTime)

	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)

	data, err := json.Marshal(&metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics body: %w", err)
	}

	logger.Info(
		ctx,
		fmt.Sprintf("metrics to send: %s", data),
	)

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(&buffer).
		Post(url)

	if err != nil {
		return fmt.Errorf("client.ReadAll: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf(
			"failed to make request: status=%s body=%s",
			resp.Status(),
			resp.Body(),
		)
	}

	logger.Info(
		ctx,
		"metrics has been sent successfully",
	)

	return nil
}

func SendMetric(
	ctx context.Context,
	baseURL string,
	metrics []*Metrics,
) error {
	url := fmt.Sprintf("%s/update/", baseURL)
	client := resty.New()

	client.
		SetHeader("Accept-Encoding", "gzip").
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime).
		SetRetryMaxWaitTime(RetryMaxWaitTime)

	for _, m := range metrics {
		var buffer bytes.Buffer
		writer := gzip.NewWriter(&buffer)

		data, err := json.Marshal(&m)
		if err != nil {
			return fmt.Errorf("failed to marshal metrics body: %w", err)
		}
		_, err = writer.Write(data)
		if err != nil {
			return fmt.Errorf("failed to compress data: %w", err)
		}

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		resp, err := client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(&buffer).
			Post(url)

		if err != nil {
			return fmt.Errorf("client.ReadAll: %w", err)
		}

		if resp.IsError() {
			return fmt.Errorf(
				"failed to make request: status=%s body=%s",
				resp.Status(),
				resp.Body(),
			)
		}

		logger.Info(
			ctx,
			fmt.Sprintf(
				"metrics has been sent successfully: widget_type=%s name=%s",
				m.MType,
				m.ID,
			),
		)
	}
	return nil
}
