package interfaces

import (
	"metrix/internal/domain"
)

type WidgetRepository struct {
	StorageHandler StorageHandler
}

func (wr *WidgetRepository) Show(
	namespace string,
	widgetType domain.WidgetType,
	name string,
) (widget domain.Widget, err error) {
	key := wr.StorageHandler.Key(namespace, string(widgetType), name)
	val, err := wr.StorageHandler.Get(key)
	if err != nil {
		return
	}

	widget.Namespace = namespace
	widget.Name = name
	widget.Type = widgetType
	widget.Value = *val

	return
}

func (wr *WidgetRepository) Update(
	namespace string,
	widgetType domain.WidgetType,
	name string,
	value uint64,
) (widget domain.Widget, err error) {
	key := wr.StorageHandler.Key(namespace, string(widgetType), name)
	val, err := wr.StorageHandler.Set(key, value)
	if err != nil {
		return
	}

	widget.Namespace = namespace
	widget.Name = name
	widget.Type = widgetType
	widget.Value = *val

	return
}

func (wr *WidgetRepository) Increment(
	namespace string,
	widgetType domain.WidgetType,
	name string,
	value uint64,
) (widget domain.Widget, err error) {
	key := wr.StorageHandler.Key(namespace, string(widgetType), name)

	prevVal, err := wr.StorageHandler.Get(key)
	if err != nil {
		return
	}

	var newValue uint64

	if prevVal == nil {
		newValue = value
	} else {
		newValue = value + *prevVal
	}

	val, err := wr.StorageHandler.Set(key, newValue)
	if err != nil {
		return
	}

	widget.Namespace = namespace
	widget.Name = name
	widget.Type = widgetType
	widget.Value = *val

	return
}
