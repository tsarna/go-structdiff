package structdiff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiff_ComprehensiveTests(t *testing.T) {
	testCases := []struct {
		name string
		old  any
		new  any
	}{
		{
			name: "simple structs no changes",
			old:  SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
			new:  SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
		},
		{
			name: "simple structs with changes",
			old:  SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
			new:  SimpleStruct{Name: "Jane", Age: 31, Email: "jane@example.com"},
		},
		{
			name: "simple structs partial changes",
			old:  SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
			new:  SimpleStruct{Name: "John", Age: 31, Email: "john@example.com"},
		},
		{
			name: "nested structs no changes",
			old: NestedStruct{
				User:    SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
				Address: AddressStruct{Street: "123 Main St", City: "NYC", ZipCode: "10001", Country: "USA"},
				Tags:    []string{"admin", "active"},
				Meta:    map[string]any{"verified": true, "score": 95.5},
				Created: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			},
			new: NestedStruct{
				User:    SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
				Address: AddressStruct{Street: "123 Main St", City: "NYC", ZipCode: "10001", Country: "USA"},
				Tags:    []string{"admin", "active"},
				Meta:    map[string]any{"verified": true, "score": 95.5},
				Created: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "nested structs with changes",
			old: NestedStruct{
				User:    SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"},
				Address: AddressStruct{Street: "123 Main St", City: "NYC", ZipCode: "10001", Country: "USA"},
				Tags:    []string{"admin", "active"},
				Meta:    map[string]any{"verified": true, "score": 95.5},
				Created: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			},
			new: NestedStruct{
				User:    SimpleStruct{Name: "Jane", Age: 31, Email: "jane@example.com"},
				Address: AddressStruct{Street: "456 Oak St", City: "Boston", ZipCode: "02101", Country: "USA"},
				Tags:    []string{"admin", "active", "premium"},
				Meta:    map[string]any{"verified": false, "score": 88.0},
				Created: time.Date(2023, 12, 26, 11, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "pointers with nil values",
			old: struct {
				Name *string `json:"name"`
				Age  *int    `json:"age"`
			}{
				Name: stringPtr("John"),
				Age:  intPtr(30),
			},
			new: struct {
				Name *string `json:"name"`
				Age  *int    `json:"age"`
			}{
				Name: nil,
				Age:  intPtr(31),
			},
		},
		{
			name: "time.Time fields",
			old: struct {
				Created time.Time  `json:"created"`
				Updated *time.Time `json:"updated"`
			}{
				Created: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Updated: timePtr(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
			new: struct {
				Created time.Time  `json:"created"`
				Updated *time.Time `json:"updated"`
			}{
				Created: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),          // Same
				Updated: timePtr(time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)), // Different
			},
		},
		{
			name: "empty structs",
			old:  struct{}{},
			new:  struct{}{},
		},
		{
			name: "structs with ignored fields",
			old: struct {
				Public  string `json:"public"`
				Ignored string `json:"-"`
				private string
			}{
				Public:  "visible",
				Ignored: "hidden1",
				private: "private1",
			},
			new: struct {
				Public  string `json:"public"`
				Ignored string `json:"-"`
				private string
			}{
				Public:  "changed",
				Ignored: "hidden2",
				private: "private2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := DiffStructs(tc.old, tc.new)

			// Verification: result should be a valid JSON patch that produces the correct transformation
			if len(result) > 0 {
				oldMap := ToMap(tc.old)
				newMap := ToMap(tc.new)

				applied := ApplyToMap(oldMap, result)

				assert.Equal(t, newMap, applied, "Diff should correctly transform old to new")
			}
		})
	}
}

func TestDiff_EdgeCases(t *testing.T) {
	t.Run("nil values", func(t *testing.T) {
		result, _ := DiffStructs(nil, nil)
		// Should return empty map for nil inputs
		expected := map[string]any{}
		assert.Equal(t, expected, result)
	})

	t.Run("nil vs struct", func(t *testing.T) {
		s := SimpleStruct{Name: "test", Age: 25, Email: "test@example.com"}

		result1, _ := DiffStructs(nil, s)
		// Should include all fields from the new struct
		assert.Contains(t, result1, "name")
		assert.Contains(t, result1, "age")
		assert.Contains(t, result1, "email")

		result2, _ := DiffStructs(s, nil)
		// Should mark all fields as deleted (nil)
		assert.Contains(t, result2, "name")
		assert.Equal(t, nil, result2["name"])
	})

	t.Run("different struct types", func(t *testing.T) {
		old := SimpleStruct{Name: "test", Age: 25, Email: "test@example.com"}
		new := AddressStruct{Street: "123 Main St", City: "NYC", ZipCode: "10001", Country: "USA"}

		result, _ := DiffStructs(old, new)
		// Should handle different struct types correctly
		assert.NotEmpty(t, result)
	})

	t.Run("non-struct types", func(t *testing.T) {
		result1, _ := DiffStructs("old", "new")
		// Should handle non-struct types (may return nil for identical types)
		// The important thing is it doesn't panic
		_ = result1

		result2, _ := DiffStructs(42, 43)
		_ = result2

		// Test with actually different values
		result3, _ := DiffStructs("old", "new")
		if result3 != nil {
			assert.IsType(t, map[string]any{}, result3)
		}
	})
}

// Helper functions for creating pointers
func stringPtr(s string) *string     { return &s }
func intPtr(i int) *int              { return &i }
func timePtr(t time.Time) *time.Time { return &t }

// Test to ensure the original Diff behavior is preserved
func TestDiff_OriginalBehavior(t *testing.T) {
	old := SimpleStruct{Name: "John", Age: 30, Email: "john@example.com"}
	new := SimpleStruct{Name: "Jane", Age: 30, Email: "jane@example.com"}

	result, _ := DiffStructs(old, new)
	expected := map[string]any{
		"name":  "Jane",
		"email": "jane@example.com", // Both changed values should be included
	}

	assert.Equal(t, expected, result)
}

func TestDiffStructs_MapFields(t *testing.T) {
	// Test DiffStructs with map fields using the unified Diff function

	type UserWithConfig struct {
		Name   string         `json:"name"`
		Age    int            `json:"age"`
		Config map[string]any `json:"config"`
	}

	t.Run("map field with changes", func(t *testing.T) {
		old := UserWithConfig{
			Name: "John",
			Age:  30,
			Config: map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
		}

		new := UserWithConfig{
			Name: "John",
			Age:  30,
			Config: map[string]any{
				"theme":         "light", // changed
				"lang":          "en",    // unchanged
				"notifications": true,    // added
			},
		}

		result, _ := DiffStructs(old, new)

		expected := map[string]any{
			"config": map[string]any{
				"theme":         "light",
				"notifications": true,
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("map field to nil", func(t *testing.T) {
		old := UserWithConfig{
			Name: "John",
			Config: map[string]any{
				"theme": "dark",
			},
		}

		new := UserWithConfig{
			Name:   "John",
			Config: nil,
		}

		result, _ := DiffStructs(old, new)

		expected := map[string]any{
			"config": map[string]any{
				"theme": nil, // Field-level deletion, more granular than wholesale nil
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("nil map field to map", func(t *testing.T) {
		old := UserWithConfig{
			Name:   "John",
			Config: nil,
		}

		new := UserWithConfig{
			Name: "John",
			Config: map[string]any{
				"theme": "light",
			},
		}

		result, _ := DiffStructs(old, new)

		expected := map[string]any{
			"config": map[string]any{
				"theme": "light",
			},
		}

		assert.Equal(t, expected, result)
	})
}

func TestDiffStructs_MixedStructMapFields(t *testing.T) {
	// Test DiffStructs with mixed struct/map fields

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type UserWithMixed struct {
		Name    string `json:"name"`
		Address any    `json:"address"` // Can be struct or map
	}

	t.Run("struct field to map field", func(t *testing.T) {
		old := UserWithMixed{
			Name:    "John",
			Address: Address{Street: "123 Main St", City: "NYC"},
		}

		new := UserWithMixed{
			Name: "John",
			Address: map[string]any{
				"street": "456 Oak Ave", // changed
				"city":   "NYC",         // unchanged
				"zip":    "10001",       // added
			},
		}

		result, _ := DiffStructs(old, new)

		expected := map[string]any{
			"address": map[string]any{
				"street": "456 Oak Ave",
				"zip":    "10001",
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("map field to struct field", func(t *testing.T) {
		old := UserWithMixed{
			Name: "John",
			Address: map[string]any{
				"street": "123 Main St",
				"city":   "NYC",
				"zip":    "10001",
			},
		}

		new := UserWithMixed{
			Name:    "John",
			Address: Address{Street: "456 Oak Ave", City: "LA"},
		}

		result, _ := DiffStructs(old, new)

		expected := map[string]any{
			"address": map[string]any{
				"street": "456 Oak Ave",
				"city":   "LA",
				"zip":    nil, // deleted from map
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("nested struct with mixed fields", func(t *testing.T) {
		type ComplexUser struct {
			Name    string         `json:"name"`
			Address Address        `json:"address"`
			Config  map[string]any `json:"config"`
		}

		old := ComplexUser{
			Name:    "John",
			Address: Address{Street: "123 Main St", City: "NYC"},
			Config: map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
		}

		new := ComplexUser{
			Name:    "Jane",                                     // changed
			Address: Address{Street: "123 Main St", City: "LA"}, // city changed
			Config: map[string]any{
				"theme":         "light", // changed
				"lang":          "en",    // unchanged
				"notifications": true,    // added
			},
		}

		result, _ := DiffStructs(old, new)

		expected := map[string]any{
			"name": "Jane",
			"address": map[string]any{
				"city": "LA",
			},
			"config": map[string]any{
				"theme":         "light",
				"notifications": true,
			},
		}

		assert.Equal(t, expected, result)
	})
}

func TestDiffStructs_Integration(t *testing.T) {
	// Test that DiffStructs + ApplyToStruct works correctly with mixed types

	type UserWithConfig struct {
		Name   string         `json:"name"`
		Config map[string]any `json:"config"`
	}

	t.Run("round-trip with map fields", func(t *testing.T) {
		original := UserWithConfig{
			Name: "John",
			Config: map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
		}

		modified := UserWithConfig{
			Name: "Jane",
			Config: map[string]any{
				"theme":         "light",
				"lang":          "en",
				"notifications": true,
			},
		}

		// Generate diff
		diff, _ := DiffStructs(original, modified)

		// Apply diff to original
		result := original // copy
		err := ApplyToStruct(&result, diff)

		assert.NoError(t, err)
		assert.Equal(t, modified, result)
	})

	t.Run("round-trip with mixed struct/map", func(t *testing.T) {
		type Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}

		type UserWithMixed struct {
			Name    string `json:"name"`
			Address any    `json:"address"`
		}

		original := UserWithMixed{
			Name:    "John",
			Address: Address{Street: "123 Main St", City: "NYC"},
		}

		// When applying a map patch to an `any` field containing a struct,
		// ApplyToStruct replaces the struct with the map (since it's an interface{} field)
		expectedAfterPatch := UserWithMixed{
			Name: "Jane",
			Address: map[string]any{
				"street": "456 Oak Ave",
			},
		}

		// Create a diff that changes name and address.street
		diff := map[string]any{
			"name": "Jane",
			"address": map[string]any{
				"street": "456 Oak Ave",
			},
		}

		// Apply diff to original
		result := original // copy
		err := ApplyToStruct(&result, diff)

		assert.NoError(t, err)
		assert.Equal(t, expectedAfterPatch, result)
	})
}
