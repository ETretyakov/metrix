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

func SendMetric(
	ctx context.Context,
	baseURL string,
	widgetType WidgetType,
	name string,
	value float64,
) error {
	url := fmt.Sprintf("%s/update/%s/%s/%f", baseURL, widgetType, name, value)

	client := resty.New()

	client.
		SetRetryCount(RetryCount).
		SetRetryWaitTime(RetryWaitTime).
		SetRetryMaxWaitTime(RetryMaxWaitTime)

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "text/plain; charset=utf-8").
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
