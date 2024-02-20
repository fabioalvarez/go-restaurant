package cmutil

import (
	"go-restaurant/internal/common/adapter/handler/http"
	"strconv"
)

// StringToUint64 is a helper function to convert a string to uint64
func StringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// ToMap is a helper function to add meta and data to a map
func ToMap(m http.Meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}
