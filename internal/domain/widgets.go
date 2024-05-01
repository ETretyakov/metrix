package domain

import "fmt"

type WidgetType string

const (
	UNKNOWN        WidgetType = "unknown"
	COUNTER_WIDGET WidgetType = "counter"
	GAUGE_WIDGET   WidgetType = "gauge"
)

func ParseWidgetType(s string) (WidgetType, error) {
	if s == string(COUNTER_WIDGET) {
		return COUNTER_WIDGET, nil
	}

	if s == string(GAUGE_WIDGET) {
		return GAUGE_WIDGET, nil
	}

	return UNKNOWN, fmt.Errorf("failed to parse widget type: %s", s)
}

type Widget struct {
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Type      WidgetType `json:"type"`
	Value     uint64     `json:"value"`
}
