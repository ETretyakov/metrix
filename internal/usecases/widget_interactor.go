package usecases

import (
	"metrix/internal/domain"
)

type WidgetInteractor struct {
	WidgetRepository WidgetRepository
}

func (wi *WidgetInteractor) Show(
	namespace string,
	widgetType domain.WidgetType,
	name string,
) (widget domain.Widget, err error) {
	widget, err = wi.WidgetRepository.Show(
		namespace,
		widgetType,
		name,
	)

	return
}

func (wi *WidgetInteractor) Update(
	namespace string,
	widgetType domain.WidgetType,
	name string,
	value uint64,
) (widget domain.Widget, err error) {
	widget, err = wi.WidgetRepository.Update(
		namespace,
		widgetType,
		name,
		value,
	)

	return
}

func (wi *WidgetInteractor) Increment(
	namespace string,
	widgetType domain.WidgetType,
	name string,
	value uint64,
) (widget domain.Widget, err error) {
	widget, err = wi.WidgetRepository.Increment(
		namespace,
		widgetType,
		name,
		value,
	)

	return
}
