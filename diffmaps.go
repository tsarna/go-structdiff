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
func DiffMaps(old, new map[string]any) map[string]any {
	if old == nil && new == nil {
		return nil
	}
	if old == nil {
		// Everything in new is an addition
		result := make(map[string]any)
		for k, v := range new {
			result[k] = v
		}
		return result
	}
	if new == nil {
		// Everything in old is a deletion
		result := make(map[string]any)
		for k := range old {
			result[k] = nil
		}
		return result
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
				diff := Diff(oldVal, newVal)
				if len(diff) > 0 {
					result[key] = diff
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

	return result
}

// valuesEqual compares two values for equality
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
		diff := DiffStructs(a, b)
		return len(diff) == 0
	}

	// For basic types, use direct comparison
	return a == b
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
