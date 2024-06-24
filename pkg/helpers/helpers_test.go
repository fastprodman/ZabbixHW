package helpers

import "testing"

// TestCompareJSONWithMap tests the CompareJSONWithMap function
func Test_CompareJSONWithMap(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		m        map[string]interface{}
		expected bool
		wantErr  bool
	}{
		{
			name:    "Equal JSON and map",
			jsonStr: `{"name": "Alice", "age": 30}`,
			m: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: true,
			wantErr:  false,
		},
		{
			name:    "Equal JSON and map with different field order",
			jsonStr: `{"age": 30, "name": "Alice"}`,
			m: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: true,
			wantErr:  false,
		},
		{
			name:    "Different JSON and map",
			jsonStr: `{"name": "Alice", "age": 31}`,
			m: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: false,
			wantErr:  false,
		},
		{
			name:    "Invalid JSON string",
			jsonStr: `{"name": "Alice", "age": 30`,
			m: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: false,
			wantErr:  true,
		},
		{
			name:    "Empty JSON string",
			jsonStr: ``,
			m: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: false,
			wantErr:  true,
		},
		{
			name:     "Empty map",
			jsonStr:  `{"name": "Alice", "age": 30}`,
			m:        map[string]interface{}{},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareJSONWithMap(tt.jsonStr, tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareJSONWithMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("CompareJSONWithMap() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// TestCompareMapsAsJSON tests the CompareMapsAsJSON function
func TestCompareMapsAsJSON(t *testing.T) {
	tests := []struct {
		name     string
		map1     map[string]interface{}
		map2     map[string]interface{}
		expected bool
		wantErr  bool
	}{
		{
			name: "Equal maps",
			map1: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			map2: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "Equal maps with different field order",
			map1: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			map2: map[string]interface{}{
				"age":  30,
				"name": "Alice",
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "Different maps",
			map1: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			map2: map[string]interface{}{
				"name": "Bob",
				"age":  30,
			},
			expected: false,
			wantErr:  false,
		},
		{
			name: "One map is empty",
			map1: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			map2:     map[string]interface{}{},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Both maps are empty",
			map1:     map[string]interface{}{},
			map2:     map[string]interface{}{},
			expected: true,
			wantErr:  false,
		},
		{
			name: "Nested maps",
			map1: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			map2: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "Different nested maps",
			map1: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			map2: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Bob",
					"age":  30,
				},
			},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareMapsAsJSON(tt.map1, tt.map2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareMapsAsJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("CompareMapsAsJSON() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// TestCompareJSONStrings tests the CompareJSONStrings function
func Test_CompareJSONStrings(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr1 string
		jsonStr2 string
		expected bool
		wantErr  bool
	}{
		{
			name:     "Equal JSON strings",
			jsonStr1: `{"name": "Alice", "age": 30}`,
			jsonStr2: `{"age": 30, "name": "Alice"}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Different JSON strings",
			jsonStr1: `{"name": "Alice", "age": 30}`,
			jsonStr2: `{"name": "Bob", "age": 30}`,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Invalid JSON string 1",
			jsonStr1: `{"name": "Alice", "age": 30`,
			jsonStr2: `{"name": "Alice", "age": 30}`,
			expected: false,
			wantErr:  true,
		},
		{
			name:     "Invalid JSON string 2",
			jsonStr1: `{"name": "Alice", "age": 30}`,
			jsonStr2: `{"name": "Alice", "age": 30`,
			expected: false,
			wantErr:  true,
		},
		{
			name:     "Empty JSON strings",
			jsonStr1: `{}`,
			jsonStr2: `{}`,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "One empty JSON string",
			jsonStr1: `{"name": "Alice", "age": 30}`,
			jsonStr2: `{}`,
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareJSONStrings(tt.jsonStr1, tt.jsonStr2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareJSONStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("CompareJSONStrings() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
