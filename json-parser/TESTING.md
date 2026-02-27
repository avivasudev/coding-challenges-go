# JSON Parser Testing Guide

This guide covers the comprehensive testing infrastructure for the JSON parser project.

## Overview

The JSON parser now includes a complete testing suite with:
- **150+ unit tests** covering tokenizer and parser functionality
- **47 integration tests** using actual JSON files from `tests/` directories
- **Performance benchmarks** for optimization tracking
- **CLI regression testing** to ensure backward compatibility
- **Automated test discovery** for easy test file management

## Test Architecture

### 1. Unit Tests (`parser/*_test.go`)

#### Core Parser Tests (`parser_test.go`)
- **ValidateJSON()** - 50+ comprehensive test cases covering all JSON structures
- **Parser functions** - Individual testing of `parseObject()`, `parseArray()`, `parseValue()`, etc.
- **Error handling** - Position tracking, error message validation
- **Edge cases** - Deep nesting, large files, Unicode content

#### Tokenizer Tests (`tokenizer_test.go`)
- **Token generation** - All token types with position tracking
- **String parsing** - All 9 escape sequences, edge cases, error conditions
- **Keyword parsing** - Case-sensitive boolean/null matching
- **Number parsing** - Integer values, edge cases
- **Whitespace handling** - Space, tab, newline, carriage return normalization
- **Invalid token detection** - Precise error reporting

#### Integration Tests (`integration_test.go`)
- **File-based testing** - Automatic discovery of JSON files in `tests/step1-5/`
- **Step-by-step validation** - Tests for each parser evolution step
- **Step 5 analysis** - Detailed reporting on missing features
- **Regression protection** - Ensures CLI compatibility maintained

### 2. Testing Infrastructure

#### TestingTokenizer Interface
Provides access to private tokenizer methods for thorough unit testing:
```go
testTokenizer := NewTestingTokenizer(input)
token := testTokenizer.ParseStringToken(startPos)
```

#### JSONError Type
Structured error reporting for better test validation:
```go
type JSONError struct {
    Message   string
    Position  int
    TokenType TokenType
}
```

## Running Tests

### Quick Start
```bash
# Run all tests with coverage
go test -v -cover ./parser/...

# Run comprehensive test suite (Go + CLI regression)
./tests/run_all_tests.sh
```

### Individual Test Suites

```bash
# Core API and parser tests
go test ./parser/ -run TestValidateJSON
go test ./parser/ -run TestParseObject
go test ./parser/ -run TestParseArray

# Tokenizer-specific tests
go test ./parser/ -run TestTokenizer
go test ./parser/ -run TestParseStringToken
go test ./parser/ -run TestParseKeywordToken

# Integration and file-based tests
go test ./parser/ -run TestJSON_FileBasedTests
go test ./parser/ -run TestStep5Analysis

# Performance benchmarks
go test -bench=. -benchmem ./parser/...

# Coverage analysis
go test -coverprofile=coverage.out ./parser/...
go tool cover -html=coverage.out -o coverage.html
```

### CLI Regression Testing

The test runner automatically validates all JSON files:
```bash
# Manual CLI testing
./json-parser tests/step1/valid.json     # Should print "Valid JSON"
./json-parser tests/step2/invalid.json   # Should print error message

# Automated CLI testing (included in run_all_tests.sh)
for file in tests/step*/*.json; do
    ./json-parser "$file"  # Results compared against expectations
done
```

## Test File Organization

### Test Data Structure
```
tests/
├── step1/          # Empty objects: {}
│   ├── valid.json      # Should pass
│   └── invalid.json    # Should fail
├── step2/          # String key-value pairs
│   ├── valid.json      # {"key": "value"}
│   ├── valid2.json     # Multiple pairs
│   ├── invalid.json    # Malformed
│   └── invalid2.json   # Additional edge case
├── step3/          # Boolean, null, numeric values
├── step4/          # Arrays and nested objects
└── step5/          # Advanced features (mixed results expected)
    ├── fail1.json      # Currently unsupported features
    ├── fail2.json      # Scientific notation, etc.
    └── ...             # 47 files total
```

### Filename Conventions
- `valid*.json` - Should parse successfully
- `invalid*.json` - Should fail parsing
- `fail*.json` - Should fail parsing
- `pass*.json` - Should parse successfully

## Current Test Coverage

### Supported Features (Steps 1-4) - 100% Pass Rate Expected
- ✅ Empty objects: `{}`
- ✅ String key-value pairs: `{"key": "value"}`
- ✅ Multiple key-value pairs: `{"k1": "v1", "k2": "v2"}`
- ✅ Boolean values: `true`, `false` (case-sensitive)
- ✅ Null values: `null`
- ✅ Positive integers: `123`, `0`, `999999`
- ✅ String escapes: `\"`, `\\`, `\/`, `\n`, `\t`, `\r`, `\b`, `\f`
- ✅ Arrays: `[]`, `["val1", "val2"]`
- ✅ Nested objects: `{"outer": {"inner": "value"}}`
- ✅ Mixed structures: `{"arr": [{"nested": true}]}`
- ✅ Deep nesting: Unlimited depth recursion
- ✅ Trailing comma detection: Properly rejected

### Partially Supported Features (Step 5) - Mixed Results Expected
- ❌ Floating-point numbers: `3.14`, `0.5`
- ❌ Negative numbers: `-42`, `-3.14`
- ❌ Scientific notation: `1e5`, `2.3e-10`
- ❌ Unicode escapes: `\u1234`
- ❌ Advanced string features: Complex Unicode handling

## Adding New Tests

### Adding Unit Tests

1. **For new parser functionality:**
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        shouldErr bool
        expected  string
    }{
        {"valid case", `{"test": true}`, false, ""},
        {"invalid case", `{"test": invalid}`, true, "expected value"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateJSON(tt.input)
            // Test logic here
        })
    }
}
```

2. **For tokenizer functionality:**
```go
func TestNewTokenType(t *testing.T) {
    testTokenizer := NewTestingTokenizer(input)
    token := testTokenizer.NextToken()

    if token.Type != EXPECTED_TYPE {
        t.Errorf("Expected %s, got %s", EXPECTED_TYPE, token.Type)
    }
}
```

### Adding Test Files

1. Create JSON file in appropriate `tests/stepX/` directory
2. Use naming convention: `valid*.json` or `invalid*.json`
3. Run tests to automatically include new file:
```bash
go test ./parser/ -run TestJSON_FileBasedTests
```

### Adding Benchmarks

```go
func BenchmarkNewFeature(b *testing.B) {
    input := "test JSON content"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = ValidateJSON(input)
    }
}
```

## Debugging Test Failures

### Common Issues

1. **Position tracking errors:**
   - Check token position calculation in tokenizer
   - Verify error messages include correct position

2. **Escape sequence handling:**
   - Test all 9 standard escape sequences individually
   - Verify unterminated string detection

3. **Case sensitivity:**
   - Ensure `true`/`false`/`null` are case-sensitive
   - `True`, `FALSE`, `NULL` should be rejected

4. **Trailing comma detection:**
   - Both objects `{"key": "value",}` and arrays `[1, 2,]` should fail

### Test Debugging Commands

```bash
# Run specific failing test with verbose output
go test -v ./parser/ -run TestSpecificFailingCase

# Run with race detection
go test -race ./parser/...

# Generate test coverage for debugging
go test -coverprofile=coverage.out ./parser/...
go tool cover -func=coverage.out | grep "functionName"

# Debug CLI behavior
./json-parser tests/step5/fail1.json  # See exact error message
echo '{"test": invalid}' | go run main.go /dev/stdin  # Test direct input
```

## Performance Expectations

### Benchmarks (on typical development machine)

```
BenchmarkNextToken                 1000000    1200 ns/op     48 B/op    2 allocs/op
BenchmarkParseStringToken           500000    2800 ns/op    128 B/op    4 allocs/op
BenchmarkValidateJSON               100000   15000 ns/op    512 B/op   12 allocs/op
BenchmarkFileBasedValidation         50000   30000 ns/op   1024 B/op   25 allocs/op
```

### Performance Targets
- Simple JSON (`{}`): < 1000 ns/op
- Complex nested JSON: < 50000 ns/op
- Memory efficiency: < 100 allocs per parse operation
- Large files (1MB+): < 10ms parsing time

## Continuous Integration

### Automated Checks
The test suite is designed for CI/CD integration:

```bash
#!/bin/bash
# CI test script
set -e

# Build check
go build -o json-parser

# Unit test with coverage requirement
go test -cover ./parser/... | grep "coverage:" | awk '{if ($3 < 90) exit 1}'

# Integration test
go test ./parser/ -run TestJSON_FileBasedTests

# CLI regression test
./tests/run_all_tests.sh
```

### Coverage Requirements
- **Minimum coverage:** 90%
- **Critical functions:** 100% coverage required
  - `ValidateJSON()`
  - `parseObject()`, `parseArray()`, `parseValue()`
  - `NextToken()`, string parsing functions

## Future Testing Improvements

### Step 5 Implementation Testing
When implementing Step 5 features:

1. **Floating-point numbers:**
```go
{"pi": 3.14159, "small": 0.001}
```

2. **Negative numbers:**
```go
{"temperature": -20, "debt": -150.50}
```

3. **Scientific notation:**
```go
{"avogadro": 6.022e23, "planck": 6.626e-34}
```

4. **Unicode escapes:**
```go
{"greeting": "Hello \u4e16\u754c"}
```

### Test Infrastructure Enhancements
- Property-based testing with random JSON generation
- Fuzzing support for edge case discovery
- Performance regression tracking
- Memory leak detection in long-running tests
- Parallel test execution optimization

---

## Summary

This comprehensive testing infrastructure provides:
- **Confidence** in parser correctness across all supported features
- **Regression protection** ensuring changes don't break existing functionality
- **Documentation** of current capabilities and limitations
- **Foundation** for implementing Step 5 and beyond
- **Performance tracking** for optimization efforts

The test suite transforms this JSON parser from a manually-tested prototype into a production-ready, well-validated component suitable for integration into larger systems.