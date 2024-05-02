package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type WidgetType string

const (
	CounterType WidgetType = "counter"
	GaugeType   WidgetType = "gauge"
	UnknownType WidgetType = "unknown"
)

func SendMetric(
	ctx context.Context,
	baseUrl string,
	widgetType WidgetType,
	name string,
	value float64,
) error {
	url := fmt.Sprintf("%s/update/%s/%s/%f", baseUrl, widgetType, name, value)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("client.GetConn: %w", err)
	}

	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client.ReadAll: %w", err)
	}

	if resp.StatusCode > 300 {
		return fmt.Errorf("client.Response: status_code=%d body=%s", resp.StatusCode, body)
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
