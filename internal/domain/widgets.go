package domain

import "fmt"

type WidgetType string

const (
	UnknownWidget WidgetType = "unknown"
	CounterWidget WidgetType = "counter"
	GaugeWidget   WidgetType = "gauge"
)

func ParseWidgetType(s string) (WidgetType, error) {
	if s == string(CounterWidget) {
		return CounterWidget, nil
	}

	if s == string(GaugeWidget) {
		return GaugeWidget, nil
	}

	return UnknownWidget, fmt.Errorf("failed to parse widget type: %s", s)
}

type Widget struct {
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Type      WidgetType `json:"type"`
	Value     float64    `json:"value"`
}

type WidgetListItem struct {
	Namespace string     `json:"namespace"`
	Name      string     `json:"name"`
	Type      WidgetType `json:"type"`
}

type WidgetList struct {
	Count int              `json:"count"`
	Items []WidgetListItem `json:"items"`
}
