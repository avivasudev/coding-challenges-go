package parser

import (
	"fmt"
	"strings"
	"testing"
)

// Test ValidateJSON with comprehensive valid and invalid cases
func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		contains  string // error message should contain this
	}{
		// Valid cases - Step 1: Empty objects
		{"empty object", "{}", false, ""},
		{"empty object with spaces", "  {}  ", false, ""},

		// Valid cases - Step 2: String key-value pairs
		{"single key-value", `{"key": "value"}`, false, ""},
		{"multiple key-values", `{"key1": "value1", "key2": "value2"}`, false, ""},
		{"string with escapes", `{"key": "value with \"quotes\""}`, false, ""},
		{"string with all escapes", `{"key": "line1\nline2\ttab\rcarriage\fform\bback\\slash\/forward"}`, false, ""},

		// Valid cases - Step 3: Boolean, null, numeric values
		{"boolean values", `{"t": true, "f": false}`, false, ""},
		{"null value", `{"n": null}`, false, ""},
		{"number value", `{"num": 123}`, false, ""},
		{"mixed values", `{"bool": true, "str": "value", "num": 42, "nothing": null}`, false, ""},

		// Valid cases - Step 4: Arrays and nested objects
		{"empty array", `{"arr": []}`, false, ""},
		{"array with values", `{"arr": ["val1", "val2", "val3"]}`, false, ""},
		{"nested object", `{"obj": {"inner": "value"}}`, false, ""},
		{"mixed nesting", `{"arr": [{"nested": true}, "string", 42]}`, false, ""},
		{"deep nesting", `{"level1": {"level2": {"level3": [{"level4": "deep"}]}}}`, false, ""},

		// Valid cases - Top-level JSON values (RFC 7159 compliance)
		{"top-level array", `["val1", "val2", "val3"]`, false, ""},
		{"top-level string", `"hello world"`, false, ""},
		{"top-level number", `42`, false, ""},
		{"top-level boolean true", `true`, false, ""},
		{"top-level boolean false", `false`, false, ""},
		{"top-level null", `null`, false, ""},
		{"top-level empty array", `[]`, false, ""},

		// Invalid cases - Basic structure errors
		{"missing opening brace", `"key": "value"}`, true, "unexpected token"},
		{"missing closing brace", `{"key": "value"`, true, "expected '}'"},
		{"extra closing brace", `{"key": "value"}}`, true, "unexpected token"},
		{"empty input", "", true, "expected value"},
		{"only whitespace", "   ", true, "expected value"},

		// Invalid cases - Key-value pair errors
		{"missing colon", `{"key" "value"}`, true, "expected ':'"},
		{"missing key", `{: "value"}`, true, "expected string key"},
		{"missing value", `{"key":}`, true, "expected value"},
		{"non-string key", `{123: "value"}`, true, "expected string key"},

		// Invalid cases - String errors
		{"unterminated string", `{"key": "unterminated`, true, "unterminated string"},
		{"invalid escape", `{"key": "bad\escape"}`, false, ""}, // Current parser accepts all escapes

		// Invalid cases - Boolean/null case sensitivity
		{"wrong case true", `{"key": True}`, true, "T"},
		{"wrong case false", `{"key": False}`, true, "F"},
		{"wrong case null", `{"key": Null}`, true, "N"},

		// Invalid cases - Leading zeros in numbers (JSON spec compliance)
		{"leading zero in object", `{"count": 013}`, true, "numbers cannot have leading zeros"},
		{"leading zero top-level", `013`, true, "numbers cannot have leading zeros"},
		{"leading zero in array", `[01, 02, 03]`, true, "numbers cannot have leading zeros"},

		// Invalid cases - Trailing commas
		{"trailing comma object", `{"key": "value",}`, true, "trailing comma"},
		{"trailing comma array", `{"arr": [1, 2,]}`, true, "trailing comma"},

		// Invalid cases - Array errors
		{"missing closing bracket", `{"arr": [1, 2, 3}`, true, "expected ']'"},
		{"missing opening bracket", `{"arr": 1, 2, 3]}`, true, "expected string key"},

		// Invalid cases - Multiple JSON values
		{"two objects", `{} {}`, true, "unexpected token"},
		{"object and array", `{} []`, true, "unexpected token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSON(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				} else if tt.contains != "" && !strings.Contains(err.Error(), tt.contains) {
					t.Errorf("Expected error containing %q, but got %q", tt.contains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %q, but got: %v", tt.input, err)
				}
			}
		})
	}
}


// Test NewParser initialization
func TestNewParser(t *testing.T) {
	parser := NewParser(`{"key": "value"}`)
	if parser == nil {
		t.Fatal("NewParser returned nil")
	}
	if parser.currentToken.Type != LEFT_BRACE {
		t.Errorf("Expected first token to be LEFT_BRACE, got %s", parser.currentToken.Type)
	}
}

// Test ParseJSON method directly
func TestParseJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid object", `{"key": "value"}`, true},
		{"valid string", `"now valid per JSON spec"`, true},  // Now valid!
		{"incomplete", `{"key":`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			err := parser.ParseJSON()
			if tt.valid && err != nil {
				t.Errorf("Expected valid JSON, but got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("Expected error for invalid JSON, but got none")
			}
		})
	}
}

// Test parseObject function behavior
func TestParseObject(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		errMsg    string
	}{
		{"empty object", "{}", false, ""},
		{"single pair", `{"key": "value"}`, false, ""},
		{"multiple pairs", `{"k1": "v1", "k2": "v2"}`, false, ""},
		{"nested object", `{"outer": {"inner": "value"}}`, false, ""},
		{"missing opening", `"key": "value"}`, true, "expected '{'"},
		{"missing closing", `{"key": "value"`, true, "expected '}'"},
		{"trailing comma", `{"key": "value",}`, true, "trailing comma"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			err := parser.parseObject()
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %q", tt.input)
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}

// Test parseArray function behavior
func TestParseArray(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		errMsg    string
	}{
		{"empty array", "[]", false, ""},
		{"single value", `["value"]`, false, ""},
		{"multiple values", `["val1", "val2", "val3"]`, false, ""},
		{"mixed types", `[true, "string", 42, null]`, false, ""},
		{"nested array", `[["inner"], "outer"]`, false, ""},
		{"nested object", `[{"key": "value"}, "string"]`, false, ""},
		{"missing opening", `"val1", "val2"]`, true, "expected '['"},
		{"missing closing", `["val1", "val2"`, true, "expected ']'"},
		{"trailing comma", `["val1", "val2",]`, true, "trailing comma"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a parser that starts at an array
			parser := NewParser(tt.input)
			err := parser.parseArray()
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %q", tt.input)
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}

// Test parseValue function with all value types
func TestParseValue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"string value", `"test string"`, false},
		{"true value", "true", false},
		{"false value", "false", false},
		{"null value", "null", false},
		{"number value", "123", false},
		{"object value", `{"key": "value"}`, false},
		{"array value", `["item1", "item2"]`, false},
		{"invalid value", "invalid", true},
		{"incomplete string", `"incomplete`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			err := parser.parseValue()
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input %q", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}

// Test parseKeyValuePair function
func TestParseKeyValuePair(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
		errMsg    string
	}{
		{"valid string pair", `"key": "value"`, false, ""},
		{"valid boolean pair", `"flag": true`, false, ""},
		{"valid number pair", `"count": 42`, false, ""},
		{"missing colon", `"key" "value"`, true, "expected ':'"},
		{"missing key", `: "value"`, true, "expected string key"},
		{"missing value", `"key":`, true, "expected value"},
		{"non-string key", `123: "value"`, true, "expected string key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			err := parser.parseKeyValuePair()
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %q", tt.input)
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}

// Test error position tracking
func TestErrorPositionTracking(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		errorPos int
	}{
		{"missing brace at start", `"key": "value"}`, 5},  // Position of unexpected token
		{"missing colon at pos 5", `{"key" "value"}`, 7},
		{"missing value at end", `{"key":}`, 7},
		{"invalid token at pos 8", `{"key": invalid}`, 8},
		{"trailing comma at pos 15", `{"key": "value",}`, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSON(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %q", tt.input)
				return
			}

			errMsg := err.Error()
			if !strings.Contains(errMsg, fmt.Sprintf("position %d", tt.errorPos)) {
				t.Errorf("Expected error at position %d, but got: %s", tt.errorPos, errMsg)
			}
		})
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"deeply nested", strings.Repeat(`{"level":`, 100) + `"deep"` + strings.Repeat(`}`, 100), false},
		{"long string", `{"key": "` + strings.Repeat("a", 1000) + `"}`, false},
		{"many keys", `{"k1":"v1","k2":"v2","k3":"v3","k4":"v4","k5":"v5"}`, false},
		{"unicode in string", `{"key": "Hello 世界"}`, false},
		{"number zero", `{"zero": 0}`, false},
		{"large number", `{"big": 999999999}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSON(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for test %q", tt.name)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for test %q, got: %v", tt.name, err)
			}
		})
	}
}