package structdiff

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ToMap converts a struct to a map[string]any representation.
// It follows JSON struct tag conventions and handles nested structures,
// slices, maps, and special types like time.Time.
//
// Rules:
// - Only exported fields are included
// - JSON tags are honored for field naming
// - Fields tagged with `json:"-"` are excluded
// - Nil pointers are omitted
// - Empty values (0, "", false, []) are included
func ToMap(v any) map[string]any {
	result := toMapValue(reflect.ValueOf(v))
	if result == nil {
		return nil
	}
	if mapResult, ok := result.(map[string]any); ok {
		return mapResult
	}
	// If result is not a map, return nil (non-struct input)
	return nil
}

func toMapValue(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	// Handle pointer: omit if nil, otherwise deref
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return toMapValue(v.Elem())
	}

	switch v.Kind() {
	case reflect.Struct:
		// Special case: time.Time
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return v.Interface()
		}

		m := make(map[string]any)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}
			name := parseName(tag, field.Name)

			fv := v.Field(i)
			if fv.Kind() == reflect.Pointer && fv.IsNil() {
				continue // omit nil pointers
			}

			val := toMapValue(fv)
			if val != nil {
				m[name] = val
			}
		}
		return m

	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			return nil
		}
		s := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			s[i] = toMapValue(v.Index(i))
		}
		return s

	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		m := make(map[string]any)
		for _, key := range v.MapKeys() {
			m[fmt.Sprint(key.Interface())] = toMapValue(v.MapIndex(key))
		}
		return m

	default:
		return v.Interface()
	}
}

func parseName(tag, fallback string) string {
	if tag == "" {
		return fallback
	}
	name := strings.Split(tag, ",")[0]
	if name == "" {
		return fallback
	}
	return name
}
