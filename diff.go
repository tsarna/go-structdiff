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
// Returns nil if both values are nil or if there are no differences.
func Diff(old, new any) map[string]any {
	// Handle nil cases
	if old == nil && new == nil {
		return nil
	}

	// Determine the types of old and new values
	oldIsStruct := isStruct(old)
	oldIsMap := isMap(old)
	newIsStruct := isStruct(new)
	newIsMap := isMap(new)

	// Handle struct-struct case
	if oldIsStruct && newIsStruct {
		return DiffStructs(old, new)
	}

	// Handle map-map case
	if oldIsMap && newIsMap {
		oldMap := old.(map[string]any)
		newMap := new.(map[string]any)
		return DiffMaps(oldMap, newMap)
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
		return DiffMaps(oldMap, newMap)
	}

	// For non-struct, non-map values, do a simple equality check
	if old == new {
		return nil
	}

	// Values are different and not structs/maps, return the new value
	// This case handles primitive types, slices, etc.
	return map[string]any{"": new}
}
