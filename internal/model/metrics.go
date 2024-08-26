// Module "model" is dedicated to storing main validation and serialisation models.
package model

import (
	"strconv"
)

// MType - the string-based type for metric types.
type MType string

// MType constants for GaugeType and CounterType.
const (
	GaugeType   MType = "gauge"
	CounterType MType = "counter"
)

// Metric - the structure for metric object serialisation.
type Metric struct {
	ID    string   `json:"id"              db:"id"      validate:"required"                       goqu:"skipupdate"`
	MType MType    `json:"type"            db:"mtype"   validate:"required,oneof=counter gauge"`
	Delta *int64   `json:"delta,omitempty" db:"delta"`
	Value *float64 `json:"value,omitempty" db:"value"`
}

// SetValue - the method that allows to encapsulate value set logic for different types.
func (m *Metric) SetValue(delta *int64, value *float64) {
	if m.MType == GaugeType {
		m.Delta = nil
		m.Value = value
		return
	}

	if m.MType == CounterType {
		m.Value = nil

		if m.Delta != nil && delta != nil {
			newDelta := *m.Delta + *delta
			m.Delta = &newDelta
			return
		}

		if delta != nil {
			m.Delta = delta
			return
		}
	}
}

// GetValue - the method that gets value for dedicated type.
func (m *Metric) GetValue() string {
	if m.MType == GaugeType {
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	}

	if m.MType == CounterType {
		return strconv.FormatInt(*m.Delta, 10)
	}

	return ""
}
