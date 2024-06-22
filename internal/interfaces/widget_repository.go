package interfaces

import (
	"errors"
	"fmt"
	"metrix/internal/domain"
	"metrix/internal/exceptions"
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
	widget.Value = val

	return
}

func (wr *WidgetRepository) Update(
	namespace string,
	widgetType domain.WidgetType,
	name string,
	value float64,
) (widget domain.Widget, err error) {
	key := wr.StorageHandler.Key(namespace, string(widgetType), name)
	val, err := wr.StorageHandler.Set(key, value)
	if err != nil {
		return
	}

	widget.Namespace = namespace
	widget.Name = name
	widget.Type = widgetType
	widget.Value = val

	return
}

func (wr *WidgetRepository) Increment(
	namespace string,
	widgetType domain.WidgetType,
	name string,
	value float64,
) (widget domain.Widget, err error) {
	key := wr.StorageHandler.Key(namespace, string(widgetType), name)

	prevVal, err := wr.StorageHandler.Get(key)
	if err != nil {
		var recordNotFound exceptions.RecordNotFoundError
		if errors.As(err, &recordNotFound) {
			prevVal = 0
		} else {
			return
		}
	}

	newValue := value + prevVal

	val, err := wr.StorageHandler.Set(key, newValue)
	if err != nil {
		return
	}

	widget.Namespace = namespace
	widget.Name = name
	widget.Type = widgetType
	widget.Value = val

	return
}

func (wr *WidgetRepository) Keys(
	namespace string,
) (keys []string, err error) {
	keys, err = wr.StorageHandler.Keys(namespace)
	if err != nil {
		err = fmt.Errorf("failed to retrieve storage keys: %w", err)
	}

	return
}
