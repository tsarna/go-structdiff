package structdiff

// Diff computes a diff/patch between two values that can be any combination of structs and maps.
// This is a unified function that automatically handles:
// - struct vs struct: uses DiffStructs
// - map vs map: uses DiffMaps
// - struct vs map: converts struct to map using ToMap, then uses DiffMaps
// - map vs struct: converts struct to map using ToMap, then uses DiffMaps
//
// The resulting map contains only the changes needed to transform old into new:
// - Keys with same values: omitted
// - Keys with different values: included with new value
// - Keys only in new: included with new value
// - Keys only in old: included with nil value (indicates deletion)
// - Nested structures: recursively diffed
//
// Returns (nil, nil) if both values are nil or if there are no differences.
// Returns (result, nil) on success, or (nil, error) if an error occurs during diffing.
func Diff(old, new any) (any, error) {
	// Handle nil cases
	if old == nil && new == nil {
		return nil, nil
	}

	// Determine the types of old and new values
	oldIsStruct := isStruct(old)
	oldIsMap := isMap(old)
	newIsStruct := isStruct(new)
	newIsMap := isMap(new)

	// Handle struct-struct case
	if oldIsStruct && newIsStruct {
		result, err := DiffStructs(old, new)
		return result, err
	}

	// Handle map-map case
	if oldIsMap && newIsMap {
		oldMap := old.(map[string]any)
		newMap := new.(map[string]any)
		result, err := DiffMaps(oldMap, newMap)
		return result, err
	}

	// Handle mixed cases: convert structs to maps and use DiffMaps
	var oldMap, newMap map[string]any

	if oldIsStruct || oldIsMap {
		if oldIsStruct {
			oldMap = ToMap(old)
		} else {
			oldMap = old.(map[string]any)
		}
	}

	if newIsStruct || newIsMap {
		if newIsStruct {
			newMap = ToMap(new)
		} else {
			newMap = new.(map[string]any)
		}
	}

	// If we have maps to compare, use DiffMaps
	if oldMap != nil || newMap != nil {
		result, err := DiffMaps(oldMap, newMap)
		return result, err
	}

	// For non-struct, non-map values, do a simple equality check
	// Use safe comparison to avoid panics
	equal := safeEqual(old, new)
	if equal {
		return nil, nil
	}

	// Values are different and not structs/maps, return the new value
	// This case handles primitive types, slices, etc.
	return new, nil
}

// safeEqual performs equality comparison without panicking on uncomparable types
func safeEqual(a, b any) bool {
	defer func() {
		// If comparison panics (e.g., comparing slices), we recover and return false
		recover()
	}()
	return a == b
}
