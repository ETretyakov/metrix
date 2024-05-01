package usecases

import "metrix/internal/domain"

type WidgetRepository interface {
	Show(
		namespace string,
		widgetType domain.WidgetType,
		name string,
	) (domain.Widget, error)

	Update(
		namespace string,
		widgetType domain.WidgetType,
		name string,
		value uint64,
	) (domain.Widget, error)

	Increment(
		namespace string,
		widgetType domain.WidgetType,
		name string,
		value uint64,
	) (domain.Widget, error)
}
