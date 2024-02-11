package utils

func Fallback[T comparable](value, nilValue, fallback T) T {
	if value == nilValue {
		return fallback
	}
	return value
}

func StringFallback(value, fallback string) string {
	return Fallback[string](value, "", fallback)
}
