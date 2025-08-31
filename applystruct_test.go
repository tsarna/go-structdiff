package structdiff

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test structures
type TestStruct struct {
	Name    string         `json:"name"`
	Age     int            `json:"age"`
	Email   string         `json:"email"`
	Active  bool           `json:"active"`
	Score   float64        `json:"score"`
	Tags    []string       `json:"tags"`
	Meta    map[string]any `json:"meta"`
	Created time.Time      `json:"created"`
}

type TestStructWithPointers struct {
	Name   *string   `json:"name"`
	Age    *int      `json:"age"`
	Active *bool     `json:"active"`
	Tags   *[]string `json:"tags"`
}

type NestedTestStruct struct {
	User    TestStruct `json:"user"`
	Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	} `json:"address"`
	Optional *TestStruct `json:"optional"`
}

func TestApplyToStruct_BasicFields(t *testing.T) {
	original := &TestStruct{
		Name:   "John",
		Age:    30,
		Email:  "john@example.com",
		Active: true,
		Score:  95.5,
	}

	patch := map[string]any{
		"name":   "Jane",
		"age":    31,
		"active": false,
		"score":  88.0,
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Equal(t, "Jane", original.Name)
	assert.Equal(t, 31, original.Age)
	assert.Equal(t, "john@example.com", original.Email) // unchanged
	assert.Equal(t, false, original.Active)
	assert.Equal(t, 88.0, original.Score)
}

func TestApplyToStruct_TypeConversions(t *testing.T) {
	original := &TestStruct{}

	patch := map[string]any{
		"name":   []byte("ByteString"),
		"age":    "42",   // string to int
		"active": "true", // string to bool
		"score":  "95.5", // string to float
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Equal(t, "ByteString", original.Name)
	assert.Equal(t, 42, original.Age)
	assert.Equal(t, true, original.Active)
	assert.Equal(t, 95.5, original.Score)
}

func TestApplyToStruct_NumericConversions(t *testing.T) {
	original := &TestStruct{}

	patch := map[string]any{
		"age":   float64(42), // float to int
		"score": int(95),     // int to float
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Equal(t, 42, original.Age)
	assert.Equal(t, 95.0, original.Score)
}

func TestApplyToStruct_SlicesAndMaps(t *testing.T) {
	original := &TestStruct{}

	patch := map[string]any{
		"tags": []any{"admin", "user", "active"},
		"meta": map[string]any{
			"verified": true,
			"score":    95.5,
			"level":    "premium",
		},
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Equal(t, []string{"admin", "user", "active"}, original.Tags)
	assert.Equal(t, map[string]any{
		"verified": true,
		"score":    95.5,
		"level":    "premium",
	}, original.Meta)
}

func TestApplyToStruct_TimeHandling(t *testing.T) {
	original := &TestStruct{}
	testTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	t.Run("time.Time value", func(t *testing.T) {
		patch := map[string]any{
			"created": testTime,
		}

		err := ApplyToStruct(original, patch)
		require.NoError(t, err)
		assert.Equal(t, testTime, original.Created)
	})

	t.Run("RFC3339 string", func(t *testing.T) {
		patch := map[string]any{
			"created": "2023-12-25T10:30:00Z",
		}

		err := ApplyToStruct(original, patch)
		require.NoError(t, err)
		assert.Equal(t, testTime, original.Created)
	})

	t.Run("invalid time string", func(t *testing.T) {
		patch := map[string]any{
			"created": "invalid-time",
		}

		err := ApplyToStruct(original, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot parse time string")
	})
}

func TestApplyToStruct_PointerFields(t *testing.T) {
	original := &TestStructWithPointers{}

	name := "John"
	age := 30
	active := true
	tags := []string{"admin", "user"}

	patch := map[string]any{
		"name":   "Jane",
		"age":    31,
		"active": false,
		"tags":   []any{"user", "premium"},
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	require.NotNil(t, original.Name)
	assert.Equal(t, "Jane", *original.Name)

	require.NotNil(t, original.Age)
	assert.Equal(t, 31, *original.Age)

	require.NotNil(t, original.Active)
	assert.Equal(t, false, *original.Active)

	require.NotNil(t, original.Tags)
	assert.Equal(t, []string{"user", "premium"}, *original.Tags)

	_ = name
	_ = age
	_ = active
	_ = tags
}

func TestApplyToStruct_NilValues(t *testing.T) {
	name := "John"
	age := 30
	tags := []string{"admin"}

	original := &TestStructWithPointers{
		Name: &name,
		Age:  &age,
		Tags: &tags,
	}

	patch := map[string]any{
		"name": nil,
		"age":  nil,
		"tags": nil,
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Nil(t, original.Name)
	assert.Nil(t, original.Age)
	assert.Nil(t, original.Tags)
}

func TestApplyToStruct_NilValueErrors(t *testing.T) {
	original := &TestStruct{
		Name: "John",
		Age:  30,
	}

	patch := map[string]any{
		"name": nil, // Cannot set string field to nil
		"age":  nil, // Cannot set int field to nil
	}

	err := ApplyToStruct(original, patch)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot set non-nillable field")
}

func TestApplyToStruct_NestedStructs(t *testing.T) {
	original := &NestedTestStruct{
		User: TestStruct{
			Name: "John",
			Age:  30,
		},
		Address: struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}{
			Street: "123 Main St",
			City:   "NYC",
		},
	}

	patch := map[string]any{
		"user": map[string]any{
			"name": "Jane",
			"age":  31,
		},
		"address": map[string]any{
			"city": "Boston",
		},
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Equal(t, "Jane", original.User.Name)
	assert.Equal(t, 31, original.User.Age)
	assert.Equal(t, "123 Main St", original.Address.Street) // unchanged
	assert.Equal(t, "Boston", original.Address.City)
}

func TestApplyToStruct_NestedPointers(t *testing.T) {
	original := &NestedTestStruct{}

	patch := map[string]any{
		"optional": map[string]any{
			"name": "Optional User",
			"age":  25,
		},
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	require.NotNil(t, original.Optional)
	assert.Equal(t, "Optional User", original.Optional.Name)
	assert.Equal(t, 25, original.Optional.Age)
}

func TestApplyToStruct_ErrorCases(t *testing.T) {
	t.Run("nil target", func(t *testing.T) {
		err := ApplyToStruct(nil, map[string]any{"name": "test"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target is nil")
	})

	t.Run("non-pointer target", func(t *testing.T) {
		target := TestStruct{}
		err := ApplyToStruct(target, map[string]any{"name": "test"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a pointer to a struct")
	})

	t.Run("pointer to non-struct", func(t *testing.T) {
		target := "string"
		err := ApplyToStruct(&target, map[string]any{"name": "test"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must point to a struct")
	})

	t.Run("field not found", func(t *testing.T) {
		target := &TestStruct{}
		err := ApplyToStruct(target, map[string]any{"nonexistent": "value"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field \"nonexistent\" not found")
	})

	t.Run("type conversion error", func(t *testing.T) {
		target := &TestStruct{}
		err := ApplyToStruct(target, map[string]any{"age": "not-a-number"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot convert string")
	})
}

func TestApplyToStruct_JSONTags(t *testing.T) {
	type TaggedStruct struct {
		PublicName  string `json:"public_name"`
		CustomField string `json:"customField"`
		Ignored     string `json:"-"`
		NoTag       string
		EmptyTag    string `json:",omitempty"`
	}

	original := &TaggedStruct{
		PublicName:  "old1",
		CustomField: "old2",
		Ignored:     "old3",
		NoTag:       "old4",
		EmptyTag:    "old5",
	}

	patch := map[string]any{
		"public_name": "new1",
		"customField": "new2",
		"NoTag":       "new4", // Uses field name when no tag
		"EmptyTag":    "new5", // Uses field name when tag is empty
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	assert.Equal(t, "new1", original.PublicName)
	assert.Equal(t, "new2", original.CustomField)
	assert.Equal(t, "old3", original.Ignored) // Should be unchanged (ignored field)
	assert.Equal(t, "new4", original.NoTag)
	assert.Equal(t, "new5", original.EmptyTag)
}

// Cross-validation test: ToMap + ApplyPatchMap should equal ApplyToStruct + ToMap
func TestApplyToStruct_CrossValidation(t *testing.T) {
	testCases := []struct {
		name     string
		original any
		patch    map[string]any
	}{
		{
			name: "basic fields",
			original: &TestStruct{
				Name:   "John",
				Age:    30,
				Active: true,
				Score:  95.5,
			},
			patch: map[string]any{
				"name":   "Jane",
				"age":    31,
				"active": false,
			},
		},
		{
			name: "with slices and maps",
			original: &TestStruct{
				Name: "John",
				Tags: []string{"admin"},
				Meta: map[string]any{"level": "basic"},
			},
			patch: map[string]any{
				"tags": []any{"admin", "premium"},
				"meta": map[string]any{"level": "premium", "verified": true},
			},
		},
		{
			name: "nested struct",
			original: &NestedTestStruct{
				User: TestStruct{Name: "John", Age: 30},
				Address: struct {
					Street string `json:"street"`
					City   string `json:"city"`
				}{Street: "123 Main St", City: "NYC"},
			},
			patch: map[string]any{
				"user": map[string]any{
					"name": "Jane",
					"age":  31,
				},
				"address": map[string]any{
					"city": "Boston",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Method 1: ToMap + ApplyPatchMap
			originalMap := ToMap(tc.original)
			result1Map := ApplyToMap(originalMap, tc.patch)

			// Method 2: ApplyToStruct + ToMap
			// We need a copy for this test
			originalCopy := copyStruct(tc.original)
			err := ApplyToStruct(originalCopy, tc.patch)
			require.NoError(t, err)
			result2Map := ToMap(originalCopy)

			// Results should be identical
			assert.Equal(t, result1Map, result2Map,
				"ToMap+ApplyPatchMap should equal ApplyToStruct+ToMap")
		})
	}
}

func TestApplyToStruct_InterfaceFields(t *testing.T) {
	type StructWithInterface struct {
		Data any `json:"data"`
	}

	original := &StructWithInterface{}

	patch := map[string]any{
		"data": map[string]any{
			"nested": "value",
			"number": 42,
		},
	}

	err := ApplyToStruct(original, patch)
	require.NoError(t, err)

	expected := map[string]any{
		"nested": "value",
		"number": 42,
	}
	assert.Equal(t, expected, original.Data)
}

func TestApplyToStruct_MapFieldPatching(t *testing.T) {
	type StructWithMapField struct {
		Settings map[string]any `json:"settings"`
		Name     string         `json:"name"`
	}

	t.Run("merge maps using ApplyToMap", func(t *testing.T) {
		original := &StructWithMapField{
			Name: "test",
			Settings: map[string]any{
				"theme":    "dark",
				"language": "en",
				"features": map[string]any{
					"notifications": true,
					"analytics":     false,
				},
			},
		}

		patch := map[string]any{
			"settings": map[string]any{
				"theme":    "light", // update existing
				"timezone": "UTC",   // add new
				"language": nil,     // delete existing
				"features": map[string]any{
					"analytics":     true, // update nested
					"beta":          true, // add nested
					"notifications": nil,  // delete nested
				},
			},
		}

		err := ApplyToStruct(original, patch)
		require.NoError(t, err)

		assert.Equal(t, "test", original.Name) // unchanged

		expected := map[string]any{
			"theme":    "light",
			"timezone": "UTC",
			// "language" deleted
			"features": map[string]any{
				"analytics": true,
				"beta":      true,
				// "notifications" deleted
			},
		}
		assert.Equal(t, expected, original.Settings)
	})

	t.Run("handle nil original map", func(t *testing.T) {
		original := &StructWithMapField{
			Name: "test",
			// Settings is nil
		}

		patch := map[string]any{
			"settings": map[string]any{
				"theme":    "light",
				"language": "en",
			},
		}

		err := ApplyToStruct(original, patch)
		require.NoError(t, err)

		expected := map[string]any{
			"theme":    "light",
			"language": "en",
		}
		assert.Equal(t, expected, original.Settings)
	})

	t.Run("replace map with nil", func(t *testing.T) {
		original := &StructWithMapField{
			Settings: map[string]any{
				"theme": "dark",
			},
		}

		patch := map[string]any{
			"settings": nil,
		}

		err := ApplyToStruct(original, patch)
		require.NoError(t, err)

		assert.Nil(t, original.Settings)
	})

	t.Run("non-map[string]any field should use original behavior", func(t *testing.T) {
		type StructWithStringMap struct {
			Config map[string]string `json:"config"`
		}

		original := &StructWithStringMap{
			Config: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		}

		patch := map[string]any{
			"config": map[string]any{
				"key1": "new_value1",
				"key3": "value3",
			},
		}

		err := ApplyToStruct(original, patch)
		require.NoError(t, err)

		// Should completely replace, not merge
		expected := map[string]string{
			"key1": "new_value1",
			"key3": "value3",
			// key2 should be gone (replaced, not merged)
		}
		assert.Equal(t, expected, original.Config)
	})
}

// Helper function to create a deep copy of a struct (for testing)
func copyStruct(src any) any {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Pointer {
		srcType := srcVal.Type().Elem()
		copy := reflect.New(srcType)
		copy.Elem().Set(srcVal.Elem())
		return copy.Interface()
	}
	return src
}
