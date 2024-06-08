package model

import "fmt"

type MType string

const (
	GaugeType   MType = "gauge"
	CounterType MType = "counter"
)

type Metric struct {
	ID    string   `json:"id"              db:"id"      validate:"required"                       goqu:"skipupdate"`
	MType MType    `json:"type"            db:"mtype"   validate:"required,oneof=counter gauge"`
	Delta *int64   `json:"delta,omitempty" db:"delta"`
	Value *float64 `json:"value,omitempty" db:"value"`
}

func (m *Metric) SetValue(delta *int64, value *float64) {
	if m.MType == GaugeType {
		m.Value = value
		return
	}

	if m.MType == CounterType {
		newDelta := *m.Delta + *delta
		m.Delta = &newDelta
		return
	}
}

func (m *Metric) GetValue() string {
	if m.MType == GaugeType {
		return fmt.Sprintf("%f", *m.Value)
	}

	if m.MType == CounterType {
		return fmt.Sprintf("%d", *m.Delta)
	}

	return ""
}
