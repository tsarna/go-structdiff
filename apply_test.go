package structdiff

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test structures for Apply function tests
type TestUser struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type TestConfig struct {
	Theme    string         `json:"theme"`
	Language string         `json:"language"`
	Features map[string]any `json:"features"`
}

func TestApply_Struct(t *testing.T) {
	t.Run("basic struct patching", func(t *testing.T) {
		target := &TestUser{
			Name:  "John",
			Age:   30,
			Email: "john@example.com",
		}

		patch := map[string]any{
			"name": "Jane",
			"age":  31,
		}

		err := Apply(target, patch)
		require.NoError(t, err)

		assert.Equal(t, "Jane", target.Name)
		assert.Equal(t, 31, target.Age)
		assert.Equal(t, "john@example.com", target.Email) // unchanged
	})

	t.Run("struct with map field", func(t *testing.T) {
		target := &TestConfig{
			Theme:    "dark",
			Language: "en",
			Features: map[string]any{
				"notifications": true,
				"analytics":     false,
			},
		}

		patch := map[string]any{
			"theme": "light",
			"features": map[string]any{
				"analytics":     true, // update existing
				"beta":          true, // add new
				"notifications": nil,  // delete existing
			},
		}

		err := Apply(target, patch)
		require.NoError(t, err)

		assert.Equal(t, "light", target.Theme)
		assert.Equal(t, "en", target.Language) // unchanged

		expected := map[string]any{
			"analytics": true,
			"beta":      true,
			// notifications should be deleted
		}
		assert.Equal(t, expected, target.Features)
	})

	t.Run("nil patch", func(t *testing.T) {
		target := &TestUser{Name: "John"}
		original := *target

		err := Apply(target, nil)
		require.NoError(t, err)

		assert.Equal(t, original, *target) // unchanged
	})
}

func TestApply_Map(t *testing.T) {
	t.Run("basic map patching", func(t *testing.T) {
		target := &map[string]any{
			"name": "John",
			"age":  30,
			"settings": map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
		}

		patch := map[string]any{
			"name": "Jane",
			"settings": map[string]any{
				"theme":         "light",
				"notifications": true,
			},
		}

		err := Apply(target, patch)
		require.NoError(t, err)

		expected := map[string]any{
			"name": "Jane",
			"age":  30, // unchanged
			"settings": map[string]any{
				"theme":         "light",
				"lang":          "en", // unchanged
				"notifications": true, // added
			},
		}
		assert.Equal(t, expected, *target)
	})

	t.Run("map with struct values", func(t *testing.T) {
		target := &map[string]any{
			"user": TestUser{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
			"count": 5,
		}

		patch := map[string]any{
			"user": map[string]any{
				"name": "Jane",
				"age":  31,
			},
			"count": 6,
		}

		err := Apply(target, patch)
		require.NoError(t, err)

		resultUser := (*target)["user"].(TestUser)
		assert.Equal(t, "Jane", resultUser.Name)
		assert.Equal(t, 31, resultUser.Age)
		assert.Equal(t, "john@example.com", resultUser.Email) // unchanged
		assert.Equal(t, 6, (*target)["count"])
	})

	t.Run("nil map target", func(t *testing.T) {
		target := &map[string]any{}
		*target = nil

		patch := map[string]any{
			"key": "value",
		}

		err := Apply(target, patch)
		require.NoError(t, err)

		expected := map[string]any{
			"key": "value",
		}
		assert.Equal(t, expected, *target)
	})

	t.Run("nil patch", func(t *testing.T) {
		original := map[string]any{"key": "value"}
		target := &original

		err := Apply(target, nil)
		require.NoError(t, err)

		assert.Equal(t, original, *target) // unchanged
	})
}

func TestApply_ErrorCases(t *testing.T) {
	t.Run("nil target", func(t *testing.T) {
		patch := map[string]any{"key": "value"}
		err := Apply(nil, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target is nil")
	})

	t.Run("non-pointer target", func(t *testing.T) {
		target := TestUser{Name: "John"}
		patch := map[string]any{"name": "Jane"}
		err := Apply(target, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target must be a pointer")
	})

	t.Run("pointer to non-struct/non-map", func(t *testing.T) {
		target := "string"
		patch := map[string]any{"key": "value"}
		err := Apply(&target, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target must point to a struct or map")
	})

	t.Run("wrong map type", func(t *testing.T) {
		target := make(map[string]string)
		patch := map[string]any{"key": "value"}
		err := Apply(&target, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "map target must be of type map[string]any")
	})

	t.Run("struct field error", func(t *testing.T) {
		target := &TestUser{Name: "John"}
		patch := map[string]any{
			"nonexistent_field": "value",
		}
		err := Apply(target, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field \"nonexistent_field\" not found")
	})
}

func TestApply_CrossCompatibility(t *testing.T) {
	// Test that Apply produces the same results as the specific functions

	t.Run("struct compatibility with ApplyToStruct", func(t *testing.T) {
		// Test data
		patch := map[string]any{
			"name": "Jane",
			"age":  31,
		}

		// Apply using the unified function
		target1 := &TestUser{Name: "John", Age: 30, Email: "john@example.com"}
		err1 := Apply(target1, patch)
		require.NoError(t, err1)

		// Apply using the specific function
		target2 := &TestUser{Name: "John", Age: 30, Email: "john@example.com"}
		err2 := ApplyToStruct(target2, patch)
		require.NoError(t, err2)

		// Results should be identical
		assert.Equal(t, *target1, *target2)
	})

	t.Run("map compatibility with ApplyToMap", func(t *testing.T) {
		// Test data
		original := map[string]any{
			"name": "John",
			"age":  30,
			"settings": map[string]any{
				"theme": "dark",
			},
		}
		patch := map[string]any{
			"name": "Jane",
			"settings": map[string]any{
				"theme":         "light",
				"notifications": true,
			},
		}

		// Apply using the unified function
		target1 := &map[string]any{}
		*target1 = copyMap(original) // Create a copy
		err1 := Apply(target1, patch)
		require.NoError(t, err1)

		// Apply using the specific function
		result2 := ApplyToMap(original, patch)

		// Results should be identical
		assert.Equal(t, *target1, result2)
	})
}

func TestApply_ComplexScenarios(t *testing.T) {
	t.Run("nested struct and map combinations", func(t *testing.T) {
		type ComplexStruct struct {
			User     TestUser       `json:"user"`
			Settings map[string]any `json:"settings"`
			Metadata map[string]any `json:"metadata"`
		}

		target := &ComplexStruct{
			User: TestUser{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
			Settings: map[string]any{
				"theme": "dark",
				"lang":  "en",
			},
			Metadata: map[string]any{
				"version": "1.0",
				"debug":   false,
			},
		}

		patch := map[string]any{
			"user": map[string]any{
				"name": "Jane",
				"age":  31,
			},
			"settings": map[string]any{
				"theme":         "light",
				"notifications": true,
				"lang":          nil, // delete
			},
			"metadata": map[string]any{
				"debug":   true,
				"feature": "enabled", // add
			},
		}

		err := Apply(target, patch)
		require.NoError(t, err)

		// Verify user struct was patched
		assert.Equal(t, "Jane", target.User.Name)
		assert.Equal(t, 31, target.User.Age)
		assert.Equal(t, "john@example.com", target.User.Email) // unchanged

		// Verify settings map was patched
		expectedSettings := map[string]any{
			"theme":         "light",
			"notifications": true,
			// lang should be deleted
		}
		assert.Equal(t, expectedSettings, target.Settings)

		// Verify metadata map was patched
		expectedMetadata := map[string]any{
			"version": "1.0", // unchanged
			"debug":   true,
			"feature": "enabled",
		}
		assert.Equal(t, expectedMetadata, target.Metadata)
	})
}

// Example function showing how to use the unified Apply function
func ExampleApply() {
	// Example 1: Applying patch to a struct
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	user := &User{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}

	patch := map[string]any{
		"name": "Jane",
		"age":  31,
	}

	Apply(user, patch)
	fmt.Printf("Struct result: %+v\n", *user)

	// Example 2: Applying patch to a map
	data := &map[string]any{
		"user": "John",
		"settings": map[string]any{
			"theme": "dark",
		},
	}

	mapPatch := map[string]any{
		"user": "Jane",
		"settings": map[string]any{
			"theme":         "light",
			"notifications": true,
		},
	}

	Apply(data, mapPatch)
	fmt.Printf("Map result: %+v\n", *data)

	// Output:
	// Struct result: {Name:Jane Age:31 Email:john@example.com}
	// Map result: map[settings:map[notifications:true theme:light] user:Jane]
}
