package utils

import (
	"net/http"
	"strconv"
)

func GetQueryParam[T comparable](r *http.Request, key string, defaultValue T) T {

	value := r.URL.Query().Get(key)

	if value == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case string:
		return any(value).(T)
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T)
		}
		return defaultValue
	case bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return any(boolValue).(T)
		}
		return defaultValue
	case float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return any(floatValue).(T)
		}
		return defaultValue
	default:
		return defaultValue
	}

}
