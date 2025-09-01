package structdiff

import "reflect"

// DiffMaps computes a diff/patch from old map to new map.
// The resulting map contains only the changes needed to transform old into new:
//
// - Keys with same values in both maps: omitted
// - Keys with different values: included with new value
// - Keys only in new: included with new value
// - Keys only in old: included with nil value (indicates deletion)
// - Nested maps: recursively diffed using DiffMaps
// - Struct values: compared using the unified Diff function for any combination of structs and maps
//
// Applying all changes in the result to the old map would produce the new map.
// Returns (result, nil) on success, or (nil, error) if an error occurs during diffing.
func DiffMaps(old, new map[string]any) (map[string]any, error) {
	if old == nil && new == nil {
		return nil, nil
	}
	if old == nil {
		// Everything in new is an addition
		result := make(map[string]any)
		for k, v := range new {
			result[k] = v
		}
		return result, nil
	}
	if new == nil {
		// Everything in old is a deletion
		result := make(map[string]any)
		for k := range old {
			result[k] = nil
		}
		return result, nil
	}

	result := make(map[string]any)

	// Track which keys we've seen in new
	seenInNew := make(map[string]bool)

	// Process all keys in new map
	for key, newVal := range new {
		seenInNew[key] = true
		oldVal, existsInOld := old[key]

		if !existsInOld {
			// Key only exists in new - include it
			result[key] = newVal
		} else if !valuesEqual(oldVal, newVal) {
			// Key exists in both but values differ
			if (isMap(oldVal) || isStruct(oldVal)) && (isMap(newVal) || isStruct(newVal)) {
				// Use unified Diff function for any combination of maps and structs
				diff, err := Diff(oldVal, newVal)
				if err != nil {
					return nil, err
				}
				if diff != nil {
					if diffMap, ok := diff.(map[string]any); ok && len(diffMap) > 0 {
						result[key] = diff
					}
				}
			} else {
				// Different values (non-map, non-struct) - include new value
				result[key] = newVal
			}
		}
		// If values are equal, omit from result
	}

	// Process keys that exist only in old (deletions)
	for key := range old {
		if !seenInNew[key] {
			result[key] = nil
		}
	}

	return result, nil
}

// valuesEqual compares two values for equality safely, handling uncomparable types
func valuesEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// For maps, we need deep comparison
	if isMap(a) && isMap(b) {
		mapA := a.(map[string]any)
		mapB := b.(map[string]any)
		return mapsEqual(mapA, mapB)
	}

	// For slices, we need deep comparison
	if isSlice(a) && isSlice(b) {
		sliceA, okA := a.([]any)
		sliceB, okB := b.([]any)
		if okA && okB {
			return slicesEqual(sliceA, sliceB)
		}
	}

	// For structs, we need deep comparison using DiffStructs
	if isStruct(a) && isStruct(b) {
		diff, err := DiffStructs(a, b)
		if err != nil {
			// If there's an error comparing structs, consider them different
			return false
		}
		return len(diff) == 0
	}

	// For basic types, use safe comparison that handles uncomparable types
	return safeEqual(a, b)
}

// mapsEqual compares two maps for deep equality
func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, exists := b[key]
		if !exists || !valuesEqual(valA, valB) {
			return false
		}
	}

	return true
}

// slicesEqual compares two slices for deep equality
func slicesEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !valuesEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

// isMap checks if a value is a map[string]any
func isMap(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

// isSlice checks if a value is a slice
func isSlice(v any) bool {
	_, ok := v.([]any)
	return ok
}

// isStruct checks if a value is a struct
func isStruct(v any) bool {
	if v == nil {
		return false
	}
	return reflect.ValueOf(v).Kind() == reflect.Struct
}
