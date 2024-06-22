package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
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

func SendMetric(
	ctx context.Context,
	baseURL string,
	widgetType WidgetType,
	name string,
	value float64,
) error {
	url := fmt.Sprintf("%s/update/", baseURL)
	client := resty.New()

	client.
		SetHeader("Accept-Encoding", "gzip").
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime).
		SetRetryMaxWaitTime(RetryMaxWaitTime)

	metrics := Metrics{
		ID:    name,
		MType: string(widgetType),
	}

	switch widgetType {
	case CounterType:
		val := int64(value)
		metrics.Delta = &val
	default:
		metrics.Value = &value
	}

	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)

	data, err := json.Marshal(&metrics)
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

	log.Info().Caller().Str("Stage", "sending-metrics").
		Msg(
			fmt.Sprintf(
				"metrics has been sent successfully: widget_type=%s name=%s value=%f",
				widgetType,
				name,
				value,
			),
		)

	return nil
}
