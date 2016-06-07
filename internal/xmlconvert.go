package internal

import (
	"reflect"
	"strconv"
	"strings"
)

func convertToXPathDouble(v interface{}) float64 {
	t := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.String:
		s := t.String()
		if s != "" {
			s = strings.TrimSpace(s)
			if s[0] != '+' {
				v, err := strconv.ParseFloat(s, 64)
				if err == nil {
					return v
				}
			}
		}
	case reflect.Float64, reflect.Float32:
		return t.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(t.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(t.Uint())
	case reflect.Bool:
		if t.Bool() {
			return 1.0
		}
		return 0.0
	}
	return 0.0
}
