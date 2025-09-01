package structdiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyPatchMap_BasicOperations(t *testing.T) {
	original := map[string]any{
		"name": "John",
		"age":  30,
		"city": "NYC",
	}

	patch := map[string]any{
		"age":   31,                 // update
		"email": "john@example.com", // add
		"city":  nil,                // delete
	}

	result := ApplyToMap(original, patch)
	expected := map[string]any{
		"name":  "John",
		"age":   31,
		"email": "john@example.com",
		// city deleted
	}

	assert.Equal(t, expected, result)

	// Original should be unchanged
	assert.Equal(t, "NYC", original["city"])
	assert.Equal(t, 30, original["age"])
}

func TestApplyPatchMap_NestedMaps(t *testing.T) {
	original := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "NYC",
				"zip":    "10001",
			},
		},
		"settings": map[string]any{
			"theme": "light",
		},
	}

	patch := map[string]any{
		"user": map[string]any{
			"age": 31, // update nested value
			"address": map[string]any{
				"street":  "456 Oak St", // update nested value
				"zip":     nil,          // delete nested value
				"country": "USA",        // add nested value
			},
			"email": "john@example.com", // add to nested map
		},
		"settings": nil, // delete entire nested map
	}

	result := ApplyToMap(original, patch)
	expected := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  31,
			"address": map[string]any{
				"street":  "456 Oak St",
				"city":    "NYC",
				"country": "USA",
				// zip deleted
			},
			"email": "john@example.com",
		},
		// settings deleted
	}

	assert.Equal(t, expected, result)
}

func TestApplyPatchMap_ReplaceMapWithValue(t *testing.T) {
	original := map[string]any{
		"config": map[string]any{
			"debug": true,
			"level": "info",
		},
	}

	patch := map[string]any{
		"config": "simple", // replace map with simple value
	}

	result := ApplyToMap(original, patch)
	expected := map[string]any{
		"config": "simple",
	}

	assert.Equal(t, expected, result)
}

func TestApplyPatchMap_ReplaceValueWithMap(t *testing.T) {
	original := map[string]any{
		"config": "simple",
	}

	patch := map[string]any{
		"config": map[string]any{
			"debug": true,
			"level": "info",
		},
	}

	result := ApplyToMap(original, patch)
	expected := map[string]any{
		"config": map[string]any{
			"debug": true,
			"level": "info",
		},
	}

	assert.Equal(t, expected, result)
}

func TestApplyPatchMap_Arrays(t *testing.T) {
	original := map[string]any{
		"tags":    []any{"go", "test"},
		"numbers": []any{1, 2, 3},
	}

	patch := map[string]any{
		"tags":    []any{"go", "test", "diff"}, // replace array
		"numbers": nil,                         // delete array
		"colors":  []any{"red", "green"},       // add array
	}

	result := ApplyToMap(original, patch)
	expected := map[string]any{
		"tags":   []any{"go", "test", "diff"},
		"colors": []any{"red", "green"},
		// numbers deleted
	}

	assert.Equal(t, expected, result)
}

func TestApplyPatchMap_EdgeCases(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		result := ApplyToMap(nil, nil)
		assert.Nil(t, result)
	})

	t.Run("original nil", func(t *testing.T) {
		patch := map[string]any{"key": "value"}
		result := ApplyToMap(nil, patch)
		expected := map[string]any{"key": "value"}
		assert.Equal(t, expected, result)
	})

	t.Run("patch nil", func(t *testing.T) {
		original := map[string]any{"key": "value"}
		result := ApplyToMap(original, nil)
		expected := map[string]any{"key": "value"}
		assert.Equal(t, expected, result)

		// Should be a copy, not the same reference
		// Modify original to verify it's a separate copy
		original["key"] = "modified"
		assert.Equal(t, "value", result["key"])
	})

	t.Run("empty patch", func(t *testing.T) {
		original := map[string]any{"key": "value"}
		patch := map[string]any{}
		result := ApplyToMap(original, patch)
		expected := map[string]any{"key": "value"}
		assert.Equal(t, expected, result)
	})
}

func TestApplyPatchMap_DeepCopy(t *testing.T) {
	original := map[string]any{
		"nested": map[string]any{
			"values": []any{1, 2, 3},
		},
	}

	patch := map[string]any{
		"new": "value",
	}

	result := ApplyToMap(original, patch)

	// Modify the original to ensure result is independent
	original["nested"].(map[string]any)["values"].([]any)[0] = 999
	original["modified"] = true

	// Result should be unchanged
	nestedValues := result["nested"].(map[string]any)["values"].([]any)
	assert.Equal(t, 1, nestedValues[0])
	assert.NotContains(t, result, "modified")
}

// Cross-validation tests: DiffMaps and ApplyPatchMap should be inverse operations
func TestDiffMapsAndApplyPatchMap_CrossValidation(t *testing.T) {
	testCases := []struct {
		name string
		old  map[string]any
		new  map[string]any
	}{
		{
			name: "simple changes",
			old: map[string]any{
				"name": "John",
				"age":  30,
				"city": "NYC",
			},
			new: map[string]any{
				"name":  "Jane",
				"age":   30,
				"email": "jane@example.com",
			},
		},
		{
			name: "nested changes",
			old: map[string]any{
				"user": map[string]any{
					"name": "John",
					"age":  30,
					"address": map[string]any{
						"street": "123 Main St",
						"city":   "NYC",
					},
				},
			},
			new: map[string]any{
				"user": map[string]any{
					"name": "John",
					"age":  31,
					"address": map[string]any{
						"street": "456 Oak St",
						"city":   "NYC",
						"zip":    "10001",
					},
				},
			},
		},
		{
			name: "complex changes",
			old: map[string]any{
				"a": 1,
				"b": map[string]any{
					"x": "old",
					"y": []any{1, 2},
				},
				"c": "delete me",
			},
			new: map[string]any{
				"a": 2,
				"b": map[string]any{
					"x": "new",
					"z": "added",
				},
				"d": "new key",
			},
		},
		{
			name: "empty string handling",
			old: map[string]any{
				"keep":   "value",
				"delete": "gone",
				"change": "",
			},
			new: map[string]any{
				"keep":   "value",
				"change": "now has value",
				"added":  "new key",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test: old -> diff -> apply -> should equal new
			diff, _ := DiffMaps(tc.old, tc.new)
			result := ApplyToMap(tc.old, diff)
			assert.Equal(t, tc.new, result, "Applying diff to old should produce new")

			// Test: new -> reverse diff -> apply -> should equal old
			reverseDiff, _ := DiffMaps(tc.new, tc.old)
			reverseResult := ApplyToMap(tc.new, reverseDiff)
			assert.Equal(t, tc.old, reverseResult, "Applying reverse diff to new should produce old")
		})
	}
}

func TestApplyPatchMap_Idempotent(t *testing.T) {
	// Applying the same patch multiple times should give the same result
	original := map[string]any{
		"name": "John",
		"age":  30,
	}

	patch := map[string]any{
		"age":   31,
		"email": "john@example.com",
	}

	result1 := ApplyToMap(original, patch)
	result2 := ApplyToMap(result1, patch)

	assert.Equal(t, result1, result2, "Applying patch should be idempotent")
}

func TestApplyPatchMap_NoSideEffects(t *testing.T) {
	// Ensure original and patch maps are not modified
	original := map[string]any{
		"nested": map[string]any{
			"value": []any{1, 2, 3},
		},
	}

	patch := map[string]any{
		"nested": map[string]any{
			"value": []any{4, 5, 6},
		},
		"new": "value",
	}

	// Take snapshots before applying
	originalSnapshot := copyMap(original)
	patchSnapshot := copyMap(patch)

	ApplyToMap(original, patch)

	// Verify no modifications
	assert.Equal(t, originalSnapshot, original, "Original should not be modified")
	assert.Equal(t, patchSnapshot, patch, "Patch should not be modified")
}

func TestApplyPatchMap_NilValuesBehavior(t *testing.T) {
	// Document the behavior: nil values in patches always mean deletion
	// This system does not support setting keys to nil values
	t.Run("nil in patch means deletion", func(t *testing.T) {
		original := map[string]any{
			"keep":   "value",
			"delete": "will be removed",
		}

		patch := map[string]any{
			"delete": nil, // This means delete the key
		}

		result := ApplyToMap(original, patch)
		expected := map[string]any{
			"keep": "value",
			// delete key should be gone
		}

		assert.Equal(t, expected, result)
		assert.NotContains(t, result, "delete", "Key with nil patch value should be deleted")
	})

	t.Run("cannot set keys to nil values", func(t *testing.T) {
		// This documents a limitation: if you want a key with nil value in the result,
		// you cannot achieve it through this patch system
		original := map[string]any{
			"existing": "value",
		}

		// If we try to add a key with nil value, it will be ignored/deleted
		patch := map[string]any{
			"new_key": nil,
		}

		result := ApplyToMap(original, patch)
		expected := map[string]any{
			"existing": "value",
			// new_key should not appear because nil means delete
		}

		assert.Equal(t, expected, result)
		assert.NotContains(t, result, "new_key", "Cannot add keys with nil values")
	})
}

func TestApplyToMap_StructPatching(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	type NestedStruct struct {
		User    TestStruct `json:"user"`
		Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		} `json:"address"`
	}

	t.Run("patch struct with map", func(t *testing.T) {
		original := map[string]any{
			"user": TestStruct{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
			"other": "value",
		}

		patch := map[string]any{
			"user": map[string]any{
				"name": "Jane",
				"age":  31,
			},
		}

		result := ApplyToMap(original, patch)

		expectedUser := TestStruct{
			Name:  "Jane",
			Age:   31,
			Email: "john@example.com", // unchanged
		}

		assert.Equal(t, expectedUser, result["user"])
		assert.Equal(t, "value", result["other"]) // unchanged
	})

	t.Run("patch nested struct", func(t *testing.T) {
		original := map[string]any{
			"data": NestedStruct{
				User: TestStruct{
					Name:  "John",
					Age:   30,
					Email: "john@example.com",
				},
				Address: struct {
					Street string `json:"street"`
					City   string `json:"city"`
				}{
					Street: "123 Main St",
					City:   "NYC",
				},
			},
		}

		patch := map[string]any{
			"data": map[string]any{
				"user": map[string]any{
					"name": "Jane",
				},
				"address": map[string]any{
					"city": "Boston",
				},
			},
		}

		result := ApplyToMap(original, patch)

		resultData := result["data"].(NestedStruct)
		assert.Equal(t, "Jane", resultData.User.Name)
		assert.Equal(t, 30, resultData.User.Age)                   // unchanged
		assert.Equal(t, "john@example.com", resultData.User.Email) // unchanged
		assert.Equal(t, "123 Main St", resultData.Address.Street)  // unchanged
		assert.Equal(t, "Boston", resultData.Address.City)
	})

	t.Run("patch struct with nil values", func(t *testing.T) {
		type StructWithPointers struct {
			Name   *string `json:"name"`
			Age    *int    `json:"age"`
			Active *bool   `json:"active"`
		}

		name := "John"
		age := 30
		active := true

		original := map[string]any{
			"user": StructWithPointers{
				Name:   &name,
				Age:    &age,
				Active: &active,
			},
		}

		patch := map[string]any{
			"user": map[string]any{
				"name":   nil, // delete
				"active": nil, // delete
			},
		}

		result := ApplyToMap(original, patch)

		resultUser := result["user"].(StructWithPointers)
		assert.Nil(t, resultUser.Name)
		assert.Equal(t, &age, resultUser.Age) // unchanged
		assert.Nil(t, resultUser.Active)
	})

	t.Run("patch fails - fallback to map replacement", func(t *testing.T) {
		type SimpleStruct struct {
			Name string `json:"name"`
		}

		original := map[string]any{
			"data": SimpleStruct{Name: "John"},
		}

		patch := map[string]any{
			"data": map[string]any{
				"nonexistent_field": "value", // This field doesn't exist in SimpleStruct
			},
		}

		result := ApplyToMap(original, patch)

		// Should fallback to replacing with the patch map since patching failed
		expected := map[string]any{
			"nonexistent_field": "value",
		}
		assert.Equal(t, expected, result["data"])
	})

	t.Run("struct in nested map", func(t *testing.T) {
		original := map[string]any{
			"level1": map[string]any{
				"level2": TestStruct{
					Name:  "John",
					Age:   30,
					Email: "john@example.com",
				},
			},
		}

		patch := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"name": "Jane",
					"age":  31,
				},
			},
		}

		result := ApplyToMap(original, patch)

		level1 := result["level1"].(map[string]any)
		level2 := level1["level2"].(TestStruct)

		assert.Equal(t, "Jane", level2.Name)
		assert.Equal(t, 31, level2.Age)
		assert.Equal(t, "john@example.com", level2.Email) // unchanged
	})

	t.Run("mixed struct and map values", func(t *testing.T) {
		original := map[string]any{
			"struct_value": TestStruct{Name: "John", Age: 30},
			"map_value": map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			"simple_value": "hello",
		}

		patch := map[string]any{
			"struct_value": map[string]any{
				"name": "Jane",
			},
			"map_value": map[string]any{
				"key1": "new_value1",
				"key3": "value3",
			},
			"simple_value": "world",
		}

		result := ApplyToMap(original, patch)

		// Struct should be patched
		structResult := result["struct_value"].(TestStruct)
		assert.Equal(t, "Jane", structResult.Name)
		assert.Equal(t, 30, structResult.Age) // unchanged

		// Map should be merged
		mapResult := result["map_value"].(map[string]any)
		assert.Equal(t, "new_value1", mapResult["key1"])
		assert.Equal(t, "value2", mapResult["key2"]) // unchanged
		assert.Equal(t, "value3", mapResult["key3"])

		// Simple value should be replaced
		assert.Equal(t, "world", result["simple_value"])
	})
}
