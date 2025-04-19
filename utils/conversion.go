package utility

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func ToInt(val string) (int, error) {
	return strconv.Atoi(val)
}

func ToFloat(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func ToBool(val string) bool {
	return strings.ToLower(val) == "true"
}

func ToTime(val string, layout string) (time.Time, error) {
	return time.Parse(layout, val)
}

func FromJSON[T any](val string) (T, error) {
	var t T
	err := json.Unmarshal([]byte(val), &t)
	return t, err
}

func ToJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}
