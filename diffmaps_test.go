package structdiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffMaps_SameValues(t *testing.T) {
	old := map[string]any{
		"name": "John",
		"age":  30,
		"city": "NYC",
	}
	new := map[string]any{
		"name": "John",
		"age":  30,
		"city": "NYC",
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_DifferentValues(t *testing.T) {
	old := map[string]any{
		"name": "John",
		"age":  30,
		"city": "NYC",
	}
	new := map[string]any{
		"name": "John",
		"age":  31,       // changed
		"city": "Boston", // changed
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"age":  31,
		"city": "Boston",
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_NewKeys(t *testing.T) {
	old := map[string]any{
		"name": "John",
		"age":  30,
	}
	new := map[string]any{
		"name":  "John",
		"age":   30,
		"city":  "NYC",              // new key
		"email": "john@example.com", // new key
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"city":  "NYC",
		"email": "john@example.com",
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_DeletedKeys(t *testing.T) {
	old := map[string]any{
		"name":  "John",
		"age":   30,
		"city":  "NYC",
		"email": "john@example.com",
	}
	new := map[string]any{
		"name": "John",
		"age":  30,
		// city and email deleted
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"city":  nil,
		"email": nil,
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_MixedChanges(t *testing.T) {
	old := map[string]any{
		"name":    "John",
		"age":     30,
		"city":    "NYC",
		"deleted": "gone",
	}
	new := map[string]any{
		"name":  "Jane",             // changed
		"age":   30,                 // same
		"city":  "Boston",           // changed
		"email": "jane@example.com", // new
		// deleted key removed
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"name":    "Jane",
		"city":    "Boston",
		"email":   "jane@example.com",
		"deleted": nil,
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_NestedMaps(t *testing.T) {
	old := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
			"address": map[string]any{
				"street": "123 Main St",
				"city":   "NYC",
			},
		},
	}
	new := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  31, // changed
			"address": map[string]any{
				"street": "456 Oak St", // changed
				"city":   "NYC",
				"zip":    "10001", // new
			},
		},
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"user": map[string]any{
			"age": 31,
			"address": map[string]any{
				"street": "456 Oak St",
				"zip":    "10001",
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_NestedMapToValue(t *testing.T) {
	old := map[string]any{
		"config": map[string]any{
			"debug": true,
			"level": "info",
		},
	}
	new := map[string]any{
		"config": "simple", // map becomes simple value
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"config": "simple",
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_ValueToNestedMap(t *testing.T) {
	old := map[string]any{
		"config": "simple",
	}
	new := map[string]any{
		"config": map[string]any{
			"debug": true,
			"level": "info",
		},
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"config": map[string]any{
			"debug": true,
			"level": "info",
		},
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_Arrays(t *testing.T) {
	old := map[string]any{
		"tags":    []any{"go", "test"},
		"numbers": []any{1, 2, 3},
	}
	new := map[string]any{
		"tags":    []any{"go", "test", "diff"}, // changed
		"numbers": []any{1, 2, 3},              // same
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"tags": []any{"go", "test", "diff"},
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_NilValues(t *testing.T) {
	old := map[string]any{
		"optional": nil,
		"required": "value",
	}
	new := map[string]any{
		"optional": "now has value",
		"required": nil, // became nil
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"optional": "now has value",
		"required": nil,
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_EdgeCases(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		result, _ := DiffMaps(nil, nil)
		assert.Nil(t, result)
	})

	t.Run("old nil", func(t *testing.T) {
		new := map[string]any{"key": "value"}
		result, _ := DiffMaps(nil, new)
		expected := map[string]any{"key": "value"}
		assert.Equal(t, expected, result)
	})

	t.Run("new nil", func(t *testing.T) {
		old := map[string]any{"key": "value"}
		result, _ := DiffMaps(old, nil)
		expected := map[string]any{"key": nil}
		assert.Equal(t, expected, result)
	})

	t.Run("both empty", func(t *testing.T) {
		old := map[string]any{}
		new := map[string]any{}
		result, _ := DiffMaps(old, new)
		expected := map[string]any{}
		assert.Equal(t, expected, result)
	})
}

func TestDiffMaps_ComplexNesting(t *testing.T) {
	old := map[string]any{
		"users": map[string]any{
			"john": map[string]any{
				"age":  30,
				"city": "NYC",
				"tags": []any{"admin", "active"},
			},
			"jane": map[string]any{
				"age":  25,
				"city": "LA",
			},
		},
		"settings": map[string]any{
			"debug": true,
		},
	}

	new := map[string]any{
		"users": map[string]any{
			"john": map[string]any{
				"age":   31, // changed
				"city":  "NYC",
				"tags":  []any{"admin", "active", "premium"}, // changed
				"email": "john@example.com",                  // new
			},
			// jane deleted
			"bob": map[string]any{ // new user
				"age":  28,
				"city": "Chicago",
			},
		},
		"settings": map[string]any{
			"debug": false,  // changed
			"theme": "dark", // new
		},
	}

	result, _ := DiffMaps(old, new)
	expected := map[string]any{
		"users": map[string]any{
			"john": map[string]any{
				"age":   31,
				"tags":  []any{"admin", "active", "premium"},
				"email": "john@example.com",
			},
			"jane": nil, // deleted
			"bob": map[string]any{
				"age":  28,
				"city": "Chicago",
			},
		},
		"settings": map[string]any{
			"debug": false,
			"theme": "dark",
		},
	}

	assert.Equal(t, expected, result)
}

func TestDiffMaps_ApplyPatch(t *testing.T) {
	// Test that applying the diff actually transforms old to new
	old := map[string]any{
		"name":    "John",
		"age":     30,
		"city":    "NYC",
		"deleted": "value",
		"nested": map[string]any{
			"a": 1,
			"b": 2,
		},
	}

	new := map[string]any{
		"name":  "Jane",
		"age":   30,
		"city":  "Boston",
		"email": "jane@example.com",
		"nested": map[string]any{
			"a": 1,
			"c": 3,
		},
	}

	diff, _ := DiffMaps(old, new)

	// Apply the diff to old to get result
	result := make(map[string]any)
	for k, v := range old {
		result[k] = v
	}

	for k, v := range diff {
		if v == nil {
			delete(result, k)
		} else if isMap(v) && isMap(result[k]) {
			// For nested maps, we'd need a recursive apply function
			// For this test, we'll just replace the entire nested map
			result[k] = v
		} else {
			result[k] = v
		}
	}

	// The result should match new (for non-nested cases)
	assert.Equal(t, "Jane", result["name"])
	assert.Equal(t, 30, result["age"])
	assert.Equal(t, "Boston", result["city"])
	assert.Equal(t, "jane@example.com", result["email"])
	assert.NotContains(t, result, "deleted")
}

func TestDiffMaps_StructValues(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
		Zip    string `json:"zip"`
	}

	t.Run("struct values with differences", func(t *testing.T) {
		old := map[string]any{
			"user": User{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
			"count": 5,
		}

		new := map[string]any{
			"user": User{
				Name:  "Jane",             // changed
				Age:   30,                 // unchanged
				Email: "jane@example.com", // changed
			},
			"count": 5, // unchanged
		}

		result, _ := DiffMaps(old, new)

		// Should only include the struct diff, not the unchanged count
		expected := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"email": "jane@example.com",
				// age is omitted because it didn't change
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("identical struct values", func(t *testing.T) {
		user := User{Name: "John", Age: 30, Email: "john@example.com"}

		old := map[string]any{
			"user":  user,
			"count": 5,
		}

		new := map[string]any{
			"user":  user, // exactly the same
			"count": 5,
		}

		result, _ := DiffMaps(old, new)

		// Should be empty since nothing changed
		assert.Empty(t, result)
	})

	t.Run("nested struct and map combination", func(t *testing.T) {
		old := map[string]any{
			"user": User{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
			"address": Address{
				Street: "123 Main St",
				City:   "NYC",
				Zip:    "10001",
			},
			"settings": map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
		}

		new := map[string]any{
			"user": User{
				Name:  "Jane",             // changed
				Age:   31,                 // changed
				Email: "john@example.com", // unchanged
			},
			"address": Address{
				Street: "456 Oak St", // changed
				City:   "NYC",        // unchanged
				Zip:    "10001",      // unchanged
			},
			"settings": map[string]any{
				"theme":         "light", // changed
				"lang":          "en",    // unchanged
				"notifications": true,    // added
			},
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": map[string]any{
				"name": "Jane",
				"age":  31,
			},
			"address": map[string]any{
				"street": "456 Oak St",
			},
			"settings": map[string]any{
				"theme":         "light",
				"notifications": true,
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("struct replaced with different type", func(t *testing.T) {
		old := map[string]any{
			"data": User{Name: "John", Age: 30},
		}

		new := map[string]any{
			"data": "simple string", // replaced struct with string
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"data": "simple string",
		}

		assert.Equal(t, expected, result)
	})

	t.Run("different struct types", func(t *testing.T) {
		old := map[string]any{
			"data": User{Name: "John", Age: 30},
		}

		new := map[string]any{
			"data": Address{Street: "123 Main St", City: "NYC"},
		}

		result, _ := DiffMaps(old, new)

		// DiffStructs falls back to map comparison for different types,
		// showing field-level differences
		expected := map[string]any{
			"data": map[string]any{
				"name":   nil,           // User field deleted
				"age":    nil,           // User field deleted
				"email":  nil,           // User field deleted
				"street": "123 Main St", // Address field added
				"city":   "NYC",         // Address field added
				"zip":    "",            // Address field added (empty)
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("struct with nil fields", func(t *testing.T) {
		type UserWithPointers struct {
			Name   *string `json:"name"`
			Age    *int    `json:"age"`
			Active *bool   `json:"active"`
		}

		name1 := "John"
		age1 := 30
		active1 := true

		name2 := "Jane"
		age2 := 31

		old := map[string]any{
			"user": UserWithPointers{
				Name:   &name1,
				Age:    &age1,
				Active: &active1,
			},
		}

		new := map[string]any{
			"user": UserWithPointers{
				Name:   &name2, // changed
				Age:    &age2,  // changed
				Active: nil,    // changed to nil
			},
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": map[string]any{
				"name":   "Jane",
				"age":    31,
				"active": nil,
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("struct deletion", func(t *testing.T) {
		old := map[string]any{
			"user": User{Name: "John", Age: 30},
			"keep": "value",
		}

		new := map[string]any{
			"keep": "value",
			// user is deleted
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": nil, // deletion marker
		}

		assert.Equal(t, expected, result)
	})

	t.Run("struct addition", func(t *testing.T) {
		old := map[string]any{
			"keep": "value",
		}

		new := map[string]any{
			"keep": "value",
			"user": User{Name: "John", Age: 30}, // added
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": User{Name: "John", Age: 30},
		}

		assert.Equal(t, expected, result)
	})
}

func TestDiffMaps_StructIntegration(t *testing.T) {
	// Test that DiffMaps + ApplyToMap works correctly with structs

	type Config struct {
		Theme    string `json:"theme"`
		Language string `json:"language"`
		Debug    bool   `json:"debug"`
	}

	original := map[string]any{
		"config": Config{
			Theme:    "dark",
			Language: "en",
			Debug:    false,
		},
		"version": "1.0",
	}

	modified := map[string]any{
		"config": Config{
			Theme:    "light", // changed
			Language: "en",    // unchanged
			Debug:    true,    // changed
		},
		"version": "1.1", // changed
	}

	// Generate diff
	diff, _ := DiffMaps(original, modified)

	// Apply diff back to original
	result := ApplyToMap(original, diff)

	// Result should match modified
	assert.Equal(t, modified, result)

	// Verify the diff structure
	expectedDiff := map[string]any{
		"config": map[string]any{
			"theme": "light",
			"debug": true,
		},
		"version": "1.1",
	}
	assert.Equal(t, expectedDiff, diff)
}

func TestDiffMaps_MixedStructMapValues(t *testing.T) {
	// Test the new functionality: struct vs map comparisons

	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	t.Run("struct to map conversion in DiffMaps", func(t *testing.T) {
		old := map[string]any{
			"user":  User{Name: "John", Age: 30, Email: "john@example.com"},
			"other": "value",
		}

		new := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"age":   31,
				"phone": "555-1234", // new field
			},
			"other": "value",
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"age":   31,
				"email": nil,        // deleted from struct
				"phone": "555-1234", // added in map
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("map to struct conversion in DiffMaps", func(t *testing.T) {
		old := map[string]any{
			"user": map[string]any{
				"name":  "John",
				"age":   30,
				"phone": "555-1234",
			},
			"other": "value",
		}

		new := map[string]any{
			"user":  User{Name: "Jane", Age: 31, Email: "jane@example.com"},
			"other": "value",
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"age":   31,
				"phone": nil,                // deleted from map
				"email": "jane@example.com", // added in struct
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("mixed struct-map with nested changes", func(t *testing.T) {
		type Config struct {
			Theme string `json:"theme"`
			Debug bool   `json:"debug"`
		}

		old := map[string]any{
			"user": User{Name: "John", Age: 30},
			"config": map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
		}

		new := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"age":   30,
				"email": "jane@example.com",
			},
			"config": Config{Theme: "light", Debug: true},
		}

		result, _ := DiffMaps(old, new)

		expected := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"email": "jane@example.com",
			},
			"config": map[string]any{
				"theme": "light",
				"lang":  nil,  // deleted from map
				"debug": true, // added in struct
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("integration test: mixed types with ApplyToMap", func(t *testing.T) {
		// Test that the diff can be applied back correctly

		original := map[string]any{
			"user": User{Name: "John", Age: 30},
		}

		modified := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"age":   31,
				"email": "jane@example.com",
			},
		}

		// Generate diff
		diff, _ := DiffMaps(original, modified)

		// Apply diff back to original
		result := ApplyToMap(original, diff)

		// The result should have the struct updated with the changes
		// ApplyToMap applies the patch to the struct, keeping it as a struct
		expected := map[string]any{
			"user": User{Name: "Jane", Age: 31, Email: "jane@example.com"},
		}
		assert.Equal(t, expected, result)
	})
}
