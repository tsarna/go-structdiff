package structdiff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToMap_BasicTypes(t *testing.T) {
	type TestStruct struct {
		String string  `json:"string"`
		Int    int     `json:"int"`
		Float  float64 `json:"float"`
		Bool   bool    `json:"bool"`
		Bytes  []byte  `json:"bytes"`
	}

	input := TestStruct{
		String: "hello",
		Int:    42,
		Float:  3.14,
		Bool:   true,
		Bytes:  []byte("world"),
	}

	result := ToMap(input)
	expected := map[string]any{
		"string": "hello",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"bytes":  []any{byte(119), byte(111), byte(114), byte(108), byte(100)}, // []byte gets converted to []any
	}

	assert.Equal(t, expected, result)
}

func TestToMap_EmptyValuesIncluded(t *testing.T) {
	type TestStruct struct {
		EmptyString string   `json:"empty_string"`
		ZeroInt     int      `json:"zero_int"`
		ZeroFloat   float64  `json:"zero_float"`
		FalseBool   bool     `json:"false_bool"`
		EmptySlice  []string `json:"empty_slice"`
	}

	input := TestStruct{
		EmptyString: "",
		ZeroInt:     0,
		ZeroFloat:   0.0,
		FalseBool:   false,
		EmptySlice:  []string{},
	}

	result := ToMap(input)
	expected := map[string]any{
		"empty_string": "",
		"zero_int":     0,
		"zero_float":   0.0,
		"false_bool":   false,
		"empty_slice":  []any{},
	}

	assert.Equal(t, expected, result)
}

func TestToMap_NilPointersOmitted(t *testing.T) {
	type TestStruct struct {
		StringPtr *string `json:"string_ptr"`
		IntPtr    *int    `json:"int_ptr"`
		BoolPtr   *bool   `json:"bool_ptr"`
	}

	input := TestStruct{
		StringPtr: nil,
		IntPtr:    nil,
		BoolPtr:   nil,
	}

	result := ToMap(input)
	expected := map[string]any{} // All nil pointers should be omitted

	assert.Equal(t, expected, result)
}

func TestToMap_NonNilPointersIncluded(t *testing.T) {
	type TestStruct struct {
		StringPtr *string `json:"string_ptr"`
		IntPtr    *int    `json:"int_ptr"`
		BoolPtr   *bool   `json:"bool_ptr"`
	}

	str := "hello"
	num := 42
	flag := true

	input := TestStruct{
		StringPtr: &str,
		IntPtr:    &num,
		BoolPtr:   &flag,
	}

	result := ToMap(input)
	expected := map[string]any{
		"string_ptr": "hello",
		"int_ptr":    42,
		"bool_ptr":   true,
	}

	assert.Equal(t, expected, result)
}

func TestToMap_IgnoresOmitempty(t *testing.T) {
	type TestStruct struct {
		EmptyWithOmit string  `json:"empty_with_omit,omitempty"`
		ZeroWithOmit  int     `json:"zero_with_omit,omitempty"`
		NilPtr        *string `json:"nil_ptr,omitempty"`
	}

	input := TestStruct{
		EmptyWithOmit: "",
		ZeroWithOmit:  0,
		NilPtr:        nil,
	}

	result := ToMap(input)
	expected := map[string]any{
		"empty_with_omit": "", // omitempty ignored, empty value included
		"zero_with_omit":  0,  // omitempty ignored, zero value included
		// nil_ptr omitted because it's a nil pointer, not because of omitempty
	}

	assert.Equal(t, expected, result)
}

func TestToMap_JSONTags(t *testing.T) {
	type TestStruct struct {
		FieldOne   string `json:"field_one"`
		FieldTwo   string `json:"customName"`
		FieldThree string `json:",omitempty"` // should use field name
		FieldFour  string `json:"-"`          // should be excluded
		FieldFive  string // no tag, should use field name
	}

	input := TestStruct{
		FieldOne:   "one",
		FieldTwo:   "two",
		FieldThree: "three",
		FieldFour:  "four",
		FieldFive:  "five",
	}

	result := ToMap(input)
	expected := map[string]any{
		"field_one":  "one",
		"customName": "two",
		"FieldThree": "three",
		// FieldFour excluded due to "-" tag
		"FieldFive": "five",
	}

	assert.Equal(t, expected, result)
}

func TestToMap_NestedStructs(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
		Count int    `json:"count"`
	}

	type Outer struct {
		Name  string `json:"name"`
		Inner Inner  `json:"inner"`
	}

	input := Outer{
		Name: "outer",
		Inner: Inner{
			Value: "inner_value",
			Count: 5,
		},
	}

	result := ToMap(input)
	expected := map[string]any{
		"name": "outer",
		"inner": map[string]any{
			"value": "inner_value",
			"count": 5,
		},
	}

	assert.Equal(t, expected, result)
}

func TestToMap_NestedStructPointers(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
	}

	type Outer struct {
		Name     string `json:"name"`
		InnerPtr *Inner `json:"inner_ptr"`
		NilPtr   *Inner `json:"nil_ptr"`
	}

	input := Outer{
		Name:     "outer",
		InnerPtr: &Inner{Value: "inner_value"},
		NilPtr:   nil,
	}

	result := ToMap(input)
	expected := map[string]any{
		"name": "outer",
		"inner_ptr": map[string]any{
			"value": "inner_value",
		},
		// nil_ptr omitted because it's a nil pointer
	}

	assert.Equal(t, expected, result)
}

func TestToMap_Slices(t *testing.T) {
	type TestStruct struct {
		StringSlice []string `json:"string_slice"`
		IntSlice    []int    `json:"int_slice"`
		EmptySlice  []string `json:"empty_slice"`
		NilSlice    []string `json:"nil_slice"`
	}

	input := TestStruct{
		StringSlice: []string{"a", "b", "c"},
		IntSlice:    []int{1, 2, 3},
		EmptySlice:  []string{},
		NilSlice:    nil,
	}

	result := ToMap(input)

	// Note: nil slice should be excluded because toMapValue returns nil for nil slices
	expected := map[string]any{
		"string_slice": []any{"a", "b", "c"},
		"int_slice":    []any{1, 2, 3},
		"empty_slice":  []any{},
		// nil_slice excluded because toMapValue returns nil for nil slices
	}

	assert.Equal(t, expected, result)
}

func TestToMap_Maps(t *testing.T) {
	type TestStruct struct {
		StringMap map[string]string `json:"string_map"`
		IntMap    map[string]int    `json:"int_map"`
		EmptyMap  map[string]string `json:"empty_map"`
		NilMap    map[string]string `json:"nil_map"`
	}

	input := TestStruct{
		StringMap: map[string]string{"key1": "value1", "key2": "value2"},
		IntMap:    map[string]int{"count": 42},
		EmptyMap:  map[string]string{},
		NilMap:    nil,
	}

	result := ToMap(input)
	expected := map[string]any{
		"string_map": map[string]any{"key1": "value1", "key2": "value2"},
		"int_map":    map[string]any{"count": 42},
		"empty_map":  map[string]any{},
		// nil_map excluded because toMapValue returns nil for nil maps
	}

	assert.Equal(t, expected, result)
}

func TestToMap_TimeTypes(t *testing.T) {
	type TestStruct struct {
		Timestamp time.Time  `json:"timestamp"`
		TimePtr   *time.Time `json:"time_ptr"`
		NilTime   *time.Time `json:"nil_time"`
	}

	testTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	input := TestStruct{
		Timestamp: testTime,
		TimePtr:   &testTime,
		NilTime:   nil,
	}

	result := ToMap(input)
	expected := map[string]any{
		"timestamp": testTime,
		"time_ptr":  testTime,
		// nil_time omitted because it's a nil pointer
	}

	assert.Equal(t, expected, result)
}

func TestToMap_UnexportedFieldsIgnored(t *testing.T) {
	type TestStruct struct {
		Public      string `json:"public"`
		private     string `json:"private"`
		alsoPrivate int
	}

	input := TestStruct{
		Public:      "visible",
		private:     "hidden",
		alsoPrivate: 42,
	}

	result := ToMap(input)
	expected := map[string]any{
		"public": "visible",
		// private fields should be ignored
	}

	assert.Equal(t, expected, result)
}

func TestToMap_ComplexNesting(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name          string         `json:"name"`
		Age           int            `json:"age"`
		Address       *Address       `json:"address"`
		Hobbies       []string       `json:"hobbies"`
		Metadata      map[string]any `json:"metadata"`
		OptionalEmail *string        `json:"optional_email"`
	}

	email := "test@example.com"
	input := Person{
		Name: "John Doe",
		Age:  30,
		Address: &Address{
			Street: "123 Main St",
			City:   "Anytown",
		},
		Hobbies: []string{"reading", "coding"},
		Metadata: map[string]any{
			"verified": true,
			"score":    95.5,
		},
		OptionalEmail: &email,
	}

	result := ToMap(input)
	expected := map[string]any{
		"name": "John Doe",
		"age":  30,
		"address": map[string]any{
			"street": "123 Main St",
			"city":   "Anytown",
		},
		"hobbies": []any{"reading", "coding"},
		"metadata": map[string]any{
			"verified": true,
			"score":    95.5,
		},
		"optional_email": "test@example.com",
	}

	assert.Equal(t, expected, result)
}

func TestToMap_ArraysVsSlices(t *testing.T) {
	type TestStruct struct {
		Array [3]string `json:"array"`
		Slice []string  `json:"slice"`
	}

	input := TestStruct{
		Array: [3]string{"a", "b", "c"},
		Slice: []string{"x", "y", "z"},
	}

	result := ToMap(input)
	expected := map[string]any{
		"array": []any{"a", "b", "c"},
		"slice": []any{"x", "y", "z"},
	}

	assert.Equal(t, expected, result)
}

func TestToMap_EdgeCases(t *testing.T) {
	t.Run("empty struct", func(t *testing.T) {
		type Empty struct{}
		result := ToMap(Empty{})
		assert.Equal(t, map[string]any{}, result)
	})

	t.Run("struct with all nil pointers", func(t *testing.T) {
		type AllNil struct {
			Ptr1 *string `json:"ptr1"`
			Ptr2 *int    `json:"ptr2"`
			Ptr3 *bool   `json:"ptr3"`
		}
		result := ToMap(AllNil{})
		assert.Equal(t, map[string]any{}, result)
	})

	t.Run("struct with all excluded fields", func(t *testing.T) {
		type AllExcluded struct {
			Field1 string `json:"-"`
			Field2 string `json:"-"`
			field3 string // unexported
		}
		input := AllExcluded{
			Field1: "excluded1",
			Field2: "excluded2",
			field3: "unexported",
		}
		result := ToMap(input)
		assert.Equal(t, map[string]any{}, result)
	})
}
