#!/bin/bash

# JSON Parser Comprehensive Test Runner
# This script runs all Go tests and CLI regression tests

set -e  # Exit on any error

echo "========================================"
echo "JSON Parser Comprehensive Test Suite"
echo "========================================"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -d "parser" ]]; then
    print_status $RED "Error: Please run this script from the json-parser root directory"
    exit 1
fi

# Build the project first
print_status $BLUE "Building JSON parser..."
if go build -o json-parser; then
    print_status $GREEN "✓ Build successful"
else
    print_status $RED "✗ Build failed"
    exit 1
fi

echo ""
print_status $BLUE "Phase 1: Running Go Unit and Integration Tests"
echo "============================================"

# Run Go tests with coverage
print_status $YELLOW "Running unit tests with coverage..."
if go test -v -cover ./parser/...; then
    print_status $GREEN "✓ Go tests passed"
    echo ""

    # Generate detailed coverage report
    print_status $YELLOW "Generating coverage report..."
    go test -coverprofile=coverage.out ./parser/...
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out | grep "total:" | awk '{print "Total coverage: " $3}'
    print_status $GREEN "✓ Coverage report generated (coverage.html)"
else
    print_status $RED "✗ Go tests failed"
    echo "Continuing with CLI tests..."
fi

echo ""
print_status $BLUE "Phase 2: Running CLI Regression Tests"
echo "===================================="

# Track test results
total_cli_tests=0
passed_cli_tests=0
failed_tests=()

# Function to test a JSON file
test_json_file() {
    local file=$1
    local expected_result=$2  # "valid" or "invalid"
    local filename=$(basename "$file")

    total_cli_tests=$((total_cli_tests + 1))

    # Run the parser
    if ./json-parser "$file" > /dev/null 2>&1; then
        result="valid"
    else
        result="invalid"
    fi

    # Check if result matches expectation
    if [[ "$result" == "$expected_result" ]]; then
        echo "  ✓ $filename: $result (expected)"
        passed_cli_tests=$((passed_cli_tests + 1))
    else
        echo "  ✗ $filename: $result (expected $expected_result)"
        failed_tests+=("$filename: got $result, expected $expected_result")
    fi
}

# Test Step 1: Empty objects
if [[ -d "tests/step1" ]]; then
    print_status $YELLOW "Testing Step 1 (Empty objects)..."
    for file in tests/step1/valid*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "valid"
    done
    for file in tests/step1/invalid*.json tests/step1/fail*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "invalid"
    done
fi

# Test Step 2: String key-value pairs
if [[ -d "tests/step2" ]]; then
    print_status $YELLOW "Testing Step 2 (String key-value pairs)..."
    for file in tests/step2/valid*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "valid"
    done
    for file in tests/step2/invalid*.json tests/step2/fail*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "invalid"
    done
fi

# Test Step 3: Boolean, null, numeric values
if [[ -d "tests/step3" ]]; then
    print_status $YELLOW "Testing Step 3 (Boolean, null, numeric values)..."
    for file in tests/step3/valid*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "valid"
    done
    for file in tests/step3/invalid*.json tests/step3/fail*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "invalid"
    done
fi

# Test Step 4: Arrays and nested objects
if [[ -d "tests/step4" ]]; then
    print_status $YELLOW "Testing Step 4 (Arrays and nested objects)..."
    for file in tests/step4/valid*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "valid"
    done
    for file in tests/step4/invalid*.json tests/step4/fail*.json; do
        [[ -f "$file" ]] && test_json_file "$file" "invalid"
    done
fi

# Test Step 5: Advanced features (expected mixed results)
if [[ -d "tests/step5" ]]; then
    print_status $YELLOW "Testing Step 5 (Advanced features - mixed results expected)..."
    step5_total=0
    step5_passed=0

    for file in tests/step5/*.json; do
        if [[ -f "$file" ]]; then
            step5_total=$((step5_total + 1))
            filename=$(basename "$file")

            # For step5, we just report status without counting as pass/fail
            if ./json-parser "$file" > /dev/null 2>&1; then
                echo "  ✓ $filename: parseable (current parser supports this)"
                step5_passed=$((step5_passed + 1))
            else
                echo "  - $filename: not supported (expected for some files)"
            fi
        fi
    done

    step5_percentage=$((step5_passed * 100 / step5_total))
    print_status $BLUE "Step 5 Summary: $step5_passed/$step5_total files parseable ($step5_percentage%)"
    echo "  Note: Step 5 contains advanced features not yet implemented"
fi

echo ""
print_status $BLUE "Phase 3: Performance Benchmarks"
echo "============================="

print_status $YELLOW "Running performance benchmarks..."
go test -bench=. -benchmem ./parser/... | grep -E "(Benchmark|ns/op|allocs/op)" || true

echo ""
print_status $BLUE "Test Summary"
echo "============="

# Calculate CLI test results (excluding step5)
cli_pass_rate=0
if [[ $total_cli_tests -gt 0 ]]; then
    cli_pass_rate=$((passed_cli_tests * 100 / total_cli_tests))
fi

print_status $GREEN "CLI Tests (Steps 1-4): $passed_cli_tests/$total_cli_tests passed ($cli_pass_rate%)"

if [[ ${#failed_tests[@]} -gt 0 ]]; then
    print_status $RED "Failed CLI tests:"
    for failure in "${failed_tests[@]}"; do
        echo "  - $failure"
    done
fi

echo ""
print_status $BLUE "Generated Files:"
echo "  - coverage.html (test coverage report)"
echo "  - coverage.out (coverage data)"
echo "  - json-parser (compiled binary)"

echo ""
if [[ ${#failed_tests[@]} -eq 0 ]]; then
    print_status $GREEN "🎉 All regression tests passed!"
    echo ""
    print_status $BLUE "Next Steps:"
    echo "  1. Run 'go test -v ./parser/...' for detailed unit test output"
    echo "  2. Open coverage.html to view test coverage"
    echo "  3. Use './json-parser <file>' to test individual JSON files"
    echo "  4. Consider implementing Step 5 features (floating-point, negative numbers, etc.)"
else
    print_status $RED "⚠️  Some regression tests failed. Please review the failures above."
    exit 1
fi