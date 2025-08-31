package structdiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test structures for Diff function tests
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

func TestDiff_StructToStruct(t *testing.T) {
	t.Run("same struct type with changes", func(t *testing.T) {
		old := User{Name: "John", Age: 30, Email: "john@example.com"}
		new := User{Name: "Jane", Age: 30, Email: "jane@example.com"}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":  "Jane",
			"email": "jane@example.com",
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("identical structs", func(t *testing.T) {
		user := User{Name: "John", Age: 30, Email: "john@example.com"}
		diff := Diff(user, user)
		assert.Empty(t, diff) // Empty map, not nil
	})

	t.Run("different struct types", func(t *testing.T) {
		old := User{Name: "John", Age: 30, Email: "john@example.com"}
		new := Address{Street: "123 Main St", City: "NYC", Zip: "10001"}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":   nil, // User fields deleted
			"age":    nil,
			"email":  nil,
			"street": "123 Main St", // Address fields added
			"city":   "NYC",
			"zip":    "10001",
		}
		assert.Equal(t, expected, diff)
	})
}

func TestDiff_MapToMap(t *testing.T) {
	t.Run("maps with changes", func(t *testing.T) {
		old := map[string]any{
			"name": "John",
			"age":  30,
			"city": "NYC",
		}
		new := map[string]any{
			"name":  "Jane",
			"age":   30,
			"email": "jane@example.com",
		}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":  "Jane",
			"email": "jane@example.com",
			"city":  nil, // deleted
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("identical maps", func(t *testing.T) {
		data := map[string]any{"name": "John", "age": 30}
		diff := Diff(data, data)
		assert.Empty(t, diff) // Empty map, not nil
	})
}

func TestDiff_StructToMap(t *testing.T) {
	t.Run("struct to map conversion", func(t *testing.T) {
		old := User{Name: "John", Age: 30, Email: "john@example.com"}
		new := map[string]any{
			"name":  "Jane",
			"age":   31,
			"phone": "555-1234", // new field
		}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":  "Jane",
			"age":   31,
			"email": nil,        // deleted from struct
			"phone": "555-1234", // added in map
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("struct to empty map", func(t *testing.T) {
		old := User{Name: "John", Age: 30}
		new := map[string]any{}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":  nil,
			"age":   nil,
			"email": nil, // ToMap includes all fields, even empty ones
		}
		assert.Equal(t, expected, diff)
	})
}

func TestDiff_MapToStruct(t *testing.T) {
	t.Run("map to struct conversion", func(t *testing.T) {
		old := map[string]any{
			"name":  "John",
			"age":   30,
			"phone": "555-1234",
		}
		new := User{Name: "Jane", Age: 31, Email: "jane@example.com"}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":  "Jane",
			"age":   31,
			"phone": nil,                // deleted from map
			"email": "jane@example.com", // added in struct
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("empty map to struct", func(t *testing.T) {
		old := map[string]any{}
		new := User{Name: "John", Age: 30}

		diff := Diff(old, new)

		expected := map[string]any{
			"name":  "John",
			"age":   30,
			"email": "", // ToMap includes all fields, even empty ones
		}
		assert.Equal(t, expected, diff)
	})
}

func TestDiff_NestedStructures(t *testing.T) {
	type NestedData struct {
		User    User           `json:"user"`
		Config  map[string]any `json:"config"`
		Version string         `json:"version"`
	}

	t.Run("nested struct with mixed changes", func(t *testing.T) {
		old := NestedData{
			User: User{Name: "John", Age: 30},
			Config: map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
			Version: "1.0",
		}

		// Convert to map for comparison
		new := map[string]any{
			"user": map[string]any{
				"name":  "Jane",             // changed
				"age":   30,                 // unchanged
				"email": "jane@example.com", // added
			},
			"config": map[string]any{
				"theme":         "light", // changed
				"lang":          "en",    // unchanged
				"notifications": true,    // added
			},
			"version": "2.0", // changed
		}

		diff := Diff(old, new)

		expected := map[string]any{
			"user": map[string]any{
				"name":  "Jane",
				"email": "jane@example.com",
			},
			"config": map[string]any{
				"theme":         "light",
				"notifications": true,
			},
			"version": "2.0",
		}
		assert.Equal(t, expected, diff)
	})
}

func TestDiff_NilAndPrimitiveCases(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		diff := Diff(nil, nil)
		assert.Nil(t, diff)
	})

	t.Run("old nil, new struct", func(t *testing.T) {
		new := User{Name: "John", Age: 30}
		diff := Diff(nil, new)

		expected := map[string]any{
			"name":  "John",
			"age":   30,
			"email": "", // ToMap includes all fields, even empty ones
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("old struct, new nil", func(t *testing.T) {
		old := User{Name: "John", Age: 30}
		diff := Diff(old, nil)

		expected := map[string]any{
			"name":  nil,
			"age":   nil,
			"email": nil, // ToMap includes all fields, even empty ones
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("old nil, new map", func(t *testing.T) {
		new := map[string]any{"name": "John"}
		diff := Diff(nil, new)

		expected := map[string]any{
			"name": "John",
		}
		assert.Equal(t, expected, diff)
	})

	t.Run("non-struct non-map values", func(t *testing.T) {
		diff := Diff("hello", "world")
		expected := map[string]any{"": "world"}
		assert.Equal(t, expected, diff)
	})

	t.Run("identical non-struct non-map values", func(t *testing.T) {
		diff := Diff("hello", "hello")
		assert.Nil(t, diff)
	})
}

func TestDiff_Integration(t *testing.T) {
	// Test that Diff + Apply works correctly for all combinations

	testCases := []struct {
		name string
		old  any
		new  any
	}{
		{
			name: "struct to struct",
			old:  User{Name: "John", Age: 30},
			new:  User{Name: "Jane", Age: 31},
		},
		{
			name: "map to map",
			old:  map[string]any{"name": "John", "age": 30},
			new:  map[string]any{"name": "Jane", "age": 31},
		},
		{
			name: "struct to map",
			old:  User{Name: "John", Age: 30},
			new:  map[string]any{"name": "Jane", "age": 31, "email": "jane@example.com"},
		},
		{
			name: "map to struct",
			old:  map[string]any{"name": "John", "age": 30, "phone": "555-1234"},
			new:  User{Name: "Jane", Age: 31, Email: "jane@example.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate diff
			diff := Diff(tc.old, tc.new)

			// Convert both to maps for comparison
			var oldMap, newMap map[string]any

			if isStruct(tc.old) {
				oldMap = ToMap(tc.old)
			} else if isMap(tc.old) {
				oldMap = tc.old.(map[string]any)
			}

			if isStruct(tc.new) {
				newMap = ToMap(tc.new)
			} else if isMap(tc.new) {
				newMap = tc.new.(map[string]any)
			}

			// Apply diff to old map
			if oldMap != nil {
				result := ApplyToMap(oldMap, diff)
				assert.Equal(t, newMap, result, "Diff + Apply should produce the new value")
			}
		})
	}
}
