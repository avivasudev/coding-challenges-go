package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// TestJSON_FileBasedTests automatically discovers and tests all JSON files
func TestJSON_FileBasedTests(t *testing.T) {
	// Define test directories and their expected behaviors
	testSteps := []struct {
		step         string
		shouldPass   bool
		description  string
	}{
		{"step1", true, "Empty objects - should all pass"},
		{"step2", true, "String key-value pairs - should all pass"},
		{"step3", true, "Boolean, null, numeric values - should all pass"},
		{"step4", true, "Arrays and nested objects - should all pass"},
		{"step5", false, "Advanced features - mixed results expected"},
	}

	baseDir := filepath.Join("..", "tests")

	for _, step := range testSteps {
		t.Run(step.step, func(t *testing.T) {
			stepDir := filepath.Join(baseDir, step.step)

			// Find all JSON files in this step directory
			files, err := filepath.Glob(filepath.Join(stepDir, "*.json"))
			if err != nil {
				t.Fatalf("Failed to find JSON files in %s: %v", stepDir, err)
			}

			if len(files) == 0 {
				t.Logf("No JSON files found in %s", stepDir)
				return
			}

			// Track results for step summary
			var passed, failed, total int
			var failures []string

			for _, file := range files {
				total++
				filename := filepath.Base(file)

				// Determine if this file should pass based on filename
				shouldPass := isValidTestFile(filename)

				// Read file content
				content, err := ioutil.ReadFile(file)
				if err != nil {
					t.Errorf("Failed to read file %s: %v", file, err)
					continue
				}

				// Test the JSON
				err = ValidateJSON(string(content))

				if shouldPass {
					if err == nil {
						passed++
						t.Logf("✓ %s: Valid JSON (expected)", filename)
					} else {
						failed++
						failMsg := fmt.Sprintf("✗ %s: Expected valid JSON, got error: %v", filename, err)
						failures = append(failures, failMsg)
						t.Errorf(failMsg)
					}
				} else {
					if err != nil {
						passed++
						t.Logf("✓ %s: Invalid JSON (expected) - %v", filename, err)
					} else {
						failed++
						failMsg := fmt.Sprintf("✗ %s: Expected invalid JSON, but parsing succeeded", filename)
						failures = append(failures, failMsg)
						t.Errorf(failMsg)
					}
				}
			}

			// Step summary
			passRate := float64(passed) / float64(total) * 100
			t.Logf("\n%s Summary: %d/%d passed (%.1f%%) - %s",
				step.step, passed, total, passRate, step.description)

			// For steps 1-4, we expect high pass rates
			if step.shouldPass && passRate < 90 {
				t.Errorf("Low pass rate for %s: %.1f%% (expected >90%%)", step.step, passRate)
			}

			// Log failures for debugging
			if len(failures) > 0 {
				t.Logf("Failures in %s:", step.step)
				for _, failure := range failures {
					t.Logf("  %s", failure)
				}
			}
		})
	}
}

// isValidTestFile determines if a file should contain valid JSON based on filename
func isValidTestFile(filename string) bool {
	filename = strings.ToLower(filename)

	// Files that should be valid
	if strings.HasPrefix(filename, "valid") {
		return true
	}
	if strings.HasPrefix(filename, "pass") {
		// Now we support top-level arrays and other JSON values!
		return true
	}

	// Files that should be invalid
	if strings.HasPrefix(filename, "invalid") {
		return false
	}
	if strings.HasPrefix(filename, "fail") {
		// Now we properly reject leading zeros
		return false
	}

	// Default assumption for ambiguous names - assume valid for now
	// We can adjust this based on actual test file analysis
	return true
}

// TestSpecificValidCases tests known valid JSON structures that should work
func TestSpecificValidCases(t *testing.T) {
	validCases := map[string]string{
		"step1_empty":         "{}",
		"step2_simple":        `{"key": "value"}`,
		"step2_multiple":      `{"key1": "value1", "key2": "value2"}`,
		"step3_boolean":       `{"flag": true, "disabled": false}`,
		"step3_null":          `{"empty": null}`,
		"step3_number":        `{"count": 42}`,
		"step3_mixed":         `{"bool": true, "str": "test", "num": 123, "nil": null}`,
		"step4_empty_array":   `{"arr": []}`,
		"step4_array_values":  `{"arr": ["val1", "val2"]}`,
		"step4_nested_object": `{"obj": {"inner": "value"}}`,
		"step4_complex":       `{"users": [{"name": "Alice", "active": true}, {"name": "Bob", "active": false}]}`,
		"step4_deep_nesting":  `{"level1": {"level2": {"level3": {"level4": "deep"}}}}`,
	}

	for name, json := range validCases {
		t.Run(name, func(t *testing.T) {
			err := ValidateJSON(json)
			if err != nil {
				t.Errorf("Expected valid JSON for %s, got error: %v", name, err)
				t.Errorf("JSON: %s", json)
			}
		})
	}
}

// TestSpecificInvalidCases tests known invalid JSON structures that should fail
func TestSpecificInvalidCases(t *testing.T) {
	invalidCases := map[string]string{
		"missing_brace":       `{"key": "value"`,
		"extra_brace":         `{"key": "value"}}`,
		"missing_colon":       `{"key" "value"}`,
		"trailing_comma_obj":  `{"key": "value",}`,
		"trailing_comma_arr":  `{"arr": [1, 2,]}`,
		"wrong_case_true":     `{"flag": True}`,
		"wrong_case_false":    `{"flag": False}`,
		"wrong_case_null":     `{"empty": Null}`,
		"unterminated_string": `{"key": "value`,
		"single_quotes":       `{'key': 'value'}`,
		"unquoted_key":        `{key: "value"}`,
		// Note: bare strings, numbers, booleans are now VALID per JSON spec
	}

	for name, json := range invalidCases {
		t.Run(name, func(t *testing.T) {
			err := ValidateJSON(json)
			if err == nil {
				t.Errorf("Expected invalid JSON for %s, but parsing succeeded", name)
				t.Errorf("JSON: %s", json)
			}
		})
	}
}

// TestStep5Analysis provides detailed analysis of Step 5 test files
func TestStep5Analysis(t *testing.T) {
	stepDir := filepath.Join("..", "tests", "step5")

	// Find all JSON files in step5 directory
	files, err := filepath.Glob(filepath.Join(stepDir, "*.json"))
	if err != nil {
		t.Fatalf("Failed to find step5 files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No step5 files found")
		return
	}

	var expectedFeatures = map[string]int{
		"floating_point":     0, // -3.14, 1.23e-4
		"negative_numbers":   0, // -42, -1.5
		"scientific_notation": 0, // 1e5, 2.3e-10
		"unicode_escapes":    0, // \u1234
		"other_failures":     0,
	}

	var totalFiles = len(files)
	var parseableFiles = 0

	t.Logf("Analyzing %d Step 5 test files for missing features...", totalFiles)

	for _, file := range files {
		filename := filepath.Base(file)

		// Read file content
		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file, err)
			continue
		}

		jsonStr := string(content)
		err = ValidateJSON(jsonStr)

		if err == nil {
			parseableFiles++
		} else {
			// Categorize the failure reason based on content
			if strings.Contains(jsonStr, "-") && (strings.Contains(jsonStr, "1") || strings.Contains(jsonStr, "2") || strings.Contains(jsonStr, "3") || strings.Contains(jsonStr, "4") || strings.Contains(jsonStr, "5") || strings.Contains(jsonStr, "6") || strings.Contains(jsonStr, "7") || strings.Contains(jsonStr, "8") || strings.Contains(jsonStr, "9") || strings.Contains(jsonStr, "0")) {
				expectedFeatures["negative_numbers"]++
			} else if strings.Contains(jsonStr, ".") || strings.Contains(jsonStr, "e") || strings.Contains(jsonStr, "E") {
				if strings.Contains(jsonStr, "e") || strings.Contains(jsonStr, "E") {
					expectedFeatures["scientific_notation"]++
				} else {
					expectedFeatures["floating_point"]++
				}
			} else if strings.Contains(jsonStr, "\\u") {
				expectedFeatures["unicode_escapes"]++
			} else {
				expectedFeatures["other_failures"]++
				t.Logf("Other failure - %s: %s (error: %v)", filename, jsonStr, err)
			}
		}
	}

	// Report analysis
	t.Logf("\nStep 5 Feature Analysis:")
	t.Logf("  Total files: %d", totalFiles)
	t.Logf("  Currently parseable: %d (%.1f%%)", parseableFiles, float64(parseableFiles)/float64(totalFiles)*100)
	t.Logf("  Expected failure reasons:")
	for feature, count := range expectedFeatures {
		if count > 0 {
			t.Logf("    %s: %d files", feature, count)
		}
	}

	// This analysis helps us understand what features we need to implement next
	if expectedFeatures["floating_point"] > 0 || expectedFeatures["negative_numbers"] > 0 || expectedFeatures["scientific_notation"] > 0 {
		t.Logf("\nRecommended Step 5 implementation priorities:")
		if expectedFeatures["negative_numbers"] > 0 {
			t.Logf("  1. Negative numbers (affects %d files)", expectedFeatures["negative_numbers"])
		}
		if expectedFeatures["floating_point"] > 0 {
			t.Logf("  2. Floating-point numbers (affects %d files)", expectedFeatures["floating_point"])
		}
		if expectedFeatures["scientific_notation"] > 0 {
			t.Logf("  3. Scientific notation (affects %d files)", expectedFeatures["scientific_notation"])
		}
		if expectedFeatures["unicode_escapes"] > 0 {
			t.Logf("  4. Unicode escape sequences (affects %d files)", expectedFeatures["unicode_escapes"])
		}
	}
}

// TestRegressionProtection ensures CLI compatibility is maintained
func TestRegressionProtection(t *testing.T) {
	// Test cases that should match the CLI behavior exactly
	knownWorkingCases := []string{
		"{}",
		`{"key": "value"}`,
		`{"key": "value", "key2": "value"}`,
		`{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`,
		`{"key": "value", "key-n": 101, "key-o": {}, "key-l": []}`,
		`{"key": "value", "key-n": 101, "key-o": {"inner key": "inner value"}, "key-l": ["list value"]}`,
	}

	for i, json := range knownWorkingCases {
		t.Run(fmt.Sprintf("known_working_%d", i+1), func(t *testing.T) {
			err := ValidateJSON(json)
			if err != nil {
				t.Errorf("Regression detected: JSON that should work failed: %s", json)
				t.Errorf("Error: %v", err)
			}
		})
	}

	knownFailingCases := []string{
		`{"key2": False}`, // Case sensitivity
		`{"key-l": ['list value']}`, // Single quotes
		`{"key": "value",}`, // Trailing comma
	}

	for i, json := range knownFailingCases {
		t.Run(fmt.Sprintf("known_failing_%d", i+1), func(t *testing.T) {
			err := ValidateJSON(json)
			if err == nil {
				t.Errorf("Regression detected: JSON that should fail succeeded: %s", json)
			}
		})
	}
}

// BenchmarkFileBasedValidation benchmarks performance on actual test files
func BenchmarkFileBasedValidation(b *testing.B) {
	// Use a representative test file
	testFile := filepath.Join("..", "tests", "step4", "valid2.json")
	content, err := ioutil.ReadFile(testFile)
	if err != nil {
		b.Skipf("Could not read test file %s: %v", testFile, err)
		return
	}

	jsonStr := string(content)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateJSON(jsonStr)
	}
}