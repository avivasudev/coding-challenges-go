# JSON Parser Project - Claude Reference

## Project Overview

A **step-by-step JSON parser implementation in Go** that incrementally adds support for increasingly complex JSON structures. The project emphasizes clean architecture, extensibility, and proper error handling.

## Current Status

### ✅ Completed Steps
- **Step 1**: Parse empty objects `{}`
- **Step 2**: Parse string key-value pairs `{"key": "value"}` and multiple pairs
- **Step 3**: Parse boolean, null, and numeric values `{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`

### 🎯 Current Capabilities
- Empty objects: `{}`
- Single key-value pairs: `{"key": "value"}`
- Multiple key-value pairs: `{"key1": "value1", "key2": "value2"}`
- **Boolean values**: `true`, `false` (case-sensitive)
- **Null values**: `null`
- **Numeric values**: positive integers like `101`
- Whitespace handling and normalization
- String escape sequence support (`\"`, `\\`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t`)
- Trailing comma detection and rejection
- Detailed error reporting with position information

## Architecture

### Two-Phase Design
1. **Tokenizer** (`parser/parser.go:51-169`):
   - Lexical analysis: breaks input into tokens
   - Supported tokens: `LEFT_BRACE`, `RIGHT_BRACE`, `STRING`, `COLON`, `COMMA`, `TRUE`, `FALSE`, `NULL`, `NUMBER`, `EOF`, `INVALID`
   - Position tracking for error reporting
   - String parsing with escape sequence handling
   - Keyword parsing with case-sensitive matching
   - Integer number parsing

2. **Recursive Descent Parser** (`parser/parser.go:172-276`):
   - Syntactic analysis: validates token sequences against JSON grammar
   - Key functions:
     - `parseObject()` - handles object structure `{ ... }`
     - `parseKeyValuePair()` - handles `"key": "value"` pairs
     - `parseValue()` - supports string, boolean, null, and number values
   - Grammar-driven approach that mirrors JSON structure

### Key Design Principles
- **Separation of Concerns**: Clean split between lexical and syntactic analysis
- **Extensibility**: Easy to add new token types and grammar rules
- **Error Handling**: Precise position and context information
- **Industry Standard**: Recursive descent approach used in production parsers

## File Structure

```
json-parser/
├── main.go                 # CLI entry point
├── parser/
│   └── parser.go          # Tokenizer + Parser implementation
├── tests/
│   ├── step1/             # Empty object tests
│   ├── step2/             # String key-value tests
│   ├── step3/             # Boolean, null, number tests
│   └── step4/             # Additional tests (TODO)
├── go.mod                 # Go module definition
├── README.md              # Project documentation
└── CLAUDE.md              # This file
```

## Building and Testing

### Build
```bash
go build -o json-parser
```

### Test Current Functionality
```bash
# Step 1: Empty objects
./json-parser tests/step1/valid.json      # Should print "Valid JSON"
./json-parser tests/step1/invalid.json    # Should print error

# Step 2: String key-value pairs
./json-parser tests/step2/valid.json      # Should print "Valid JSON"
./json-parser tests/step2/valid2.json     # Should print "Valid JSON"
./json-parser tests/step2/invalid.json    # Should print error
./json-parser tests/step2/invalid2.json   # Should print error

# Step 3: Boolean, null, and numeric values
./json-parser tests/step3/valid.json      # Should print "Valid JSON"
./json-parser tests/step3/invalid.json    # Should print error
```

### Test Files Content
- `tests/step1/valid.json`: `{}`
- `tests/step2/valid.json`: `{"key": "value"}`
- `tests/step2/valid2.json`: `{"key": "value", "key2": "value"}`
- `tests/step3/valid.json`: `{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`
- `tests/step3/invalid.json`: `{"key2": False}` (case sensitivity test - "False" should be rejected)

## ✅ Step 3 Implementation (Completed)

**Successfully implemented support for:**
1. **Boolean values**: `true`, `false` (case-sensitive matching)
2. **Null values**: `null`
3. **Number values**: positive integers like `101`

**Implementation details:**
- Added new token types: `TRUE`, `FALSE`, `NULL`, `NUMBER`
- Implemented `parseKeywordToken()` for case-sensitive keyword matching
- Implemented `parseNumberToken()` for integer parsing
- Updated `NextToken()` character dispatch for letters and digits
- Extended `parseValue()` to accept multiple value types
- Maintained position tracking for precise error reporting

## Next Steps (Step 4)

Based on `tests/step4/valid.json`, the next iteration needs to support:

1. **Nested objects**: `{"key-o": {"inner key": "inner value"}}`
2. **Arrays**: `{"key-l": ["list value"]}`
3. **Mixed structures**: Objects and arrays as values

### Architecture Benefits for Extension
- Token-based approach makes adding new structural elements straightforward
- Recursive descent structure naturally accommodates nested constructs
- Existing value parsing logic can be reused for array elements and nested objects
- Error reporting system ready for complex nested structures

## Notes for Future Development

- The recursive descent approach will scale well for nested objects and arrays
- Consider adding an AST (Abstract Syntax Tree) representation for step 4+
- Current parser validates but doesn't build a data structure (validation-only)
- Position tracking enables precise error messages for debugging
- Modular design allows independent testing of tokenizer vs parser logic