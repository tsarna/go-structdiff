package structdiff

import (
	"testing"
)

type BenchmarkStruct struct {
	Name     string         `json:"name"`
	Age      int            `json:"age"`
	Active   bool           `json:"active"`
	Score    float64        `json:"score"`
	Tags     []string       `json:"tags"`
	Meta     map[string]any `json:"meta"`
	Optional *string        `json:"optional"`
	Nested   struct {
		Value string `json:"value"`
		Count int    `json:"count"`
	} `json:"nested"`
}

func BenchmarkApplyToStruct(b *testing.B) {
	original := &BenchmarkStruct{
		Name:   "John",
		Age:    30,
		Active: true,
		Score:  95.5,
		Tags:   []string{"admin", "user"},
		Meta:   map[string]any{"level": "premium", "verified": true},
		Nested: struct {
			Value string `json:"value"`
			Count int    `json:"count"`
		}{Value: "nested", Count: 42},
	}

	patch := map[string]any{
		"name":   "Jane",
		"age":    31,
		"active": false,
		"tags":   []any{"admin", "premium"},
		"meta": map[string]any{
			"level":    "platinum",
			"verified": true,
			"score":    98.5,
		},
		"nested": map[string]any{
			"value": "updated",
			"count": 100,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		target := copyBenchmarkStruct(original)
		_ = ApplyToStruct(target, patch)
	}
}

func BenchmarkToMapPlusApplyPatchMap(b *testing.B) {
	original := &BenchmarkStruct{
		Name:   "John",
		Age:    30,
		Active: true,
		Score:  95.5,
		Tags:   []string{"admin", "user"},
		Meta:   map[string]any{"level": "premium", "verified": true},
		Nested: struct {
			Value string `json:"value"`
			Count int    `json:"count"`
		}{Value: "nested", Count: 42},
	}

	patch := map[string]any{
		"name":   "Jane",
		"age":    31,
		"active": false,
		"tags":   []any{"admin", "premium"},
		"meta": map[string]any{
			"level":    "platinum",
			"verified": true,
			"score":    98.5,
		},
		"nested": map[string]any{
			"value": "updated",
			"count": 100,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		originalMap := ToMap(original)
		_ = ApplyToMap(originalMap, patch)
	}
}

func BenchmarkApplyToStruct_Allocs(b *testing.B) {
	original := &BenchmarkStruct{
		Name:   "John",
		Age:    30,
		Active: true,
		Score:  95.5,
		Tags:   []string{"admin", "user"},
		Meta:   map[string]any{"level": "premium", "verified": true},
		Nested: struct {
			Value string `json:"value"`
			Count int    `json:"count"`
		}{Value: "nested", Count: 42},
	}

	patch := map[string]any{
		"name":   "Jane",
		"age":    31,
		"active": false,
		"nested": map[string]any{
			"value": "updated",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		target := copyBenchmarkStruct(original)
		_ = ApplyToStruct(target, patch)
	}
}

func copyBenchmarkStruct(src *BenchmarkStruct) *BenchmarkStruct {
	optional := ""
	if src.Optional != nil {
		optional = *src.Optional
	}

	return &BenchmarkStruct{
		Name:     src.Name,
		Age:      src.Age,
		Active:   src.Active,
		Score:    src.Score,
		Tags:     append([]string(nil), src.Tags...),
		Meta:     copyMapBench(src.Meta),
		Optional: &optional,
		Nested: struct {
			Value string `json:"value"`
			Count int    `json:"count"`
		}{
			Value: src.Nested.Value,
			Count: src.Nested.Count,
		},
	}
}

func copyMapBench(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
