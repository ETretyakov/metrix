package client

import (
	"context"
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

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
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
