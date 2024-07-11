package monitoring

import (
	"context"
	"fmt"
	"metrix/pkg/logger"
	"reflect"
)

func fieldToFloat64(
	ctx context.Context,
	field reflect.Value,
) float64 {
	switch field.Kind() {
	case reflect.Uint32:
		if val, ok := field.Interface().(uint32); !ok {
			logger.Warn(
				ctx,
				"failed to assert filetype Uint32",
			)
		} else {
			return float64(val)
		}
	case reflect.Uint64:
		if val, ok := field.Interface().(uint64); !ok {
			logger.Warn(
				ctx,
				"failed to assert filetype Uint64",
			)
		} else {
			return float64(val)
		}
	case reflect.Float64:
		if val, ok := field.Interface().(float64); !ok {
			logger.Warn(
				ctx,
				"failed to assert filetype Float64",
			)
		} else {
			return float64(val)
		}
	default:
		logger.Warn(
			ctx,
			fmt.Sprintf(
				"unsupported metric field type: %s",
				field.Kind(),
			),
		)
	}

	return 0
}
