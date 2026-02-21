# JSON Parser

## Overview

A step-by-step implementation of a JSON parser in Go, built incrementally to support increasingly complex JSON structures. The project emphasizes clean architecture, extensibility, and proper error handling.

**Current Status:** ✅ **Step 4 Complete** - Full support for nested objects, arrays, and all JSON primitive types with unlimited nesting depth.

## Step 1: Empty Objects

**Objective:** Parse empty JSON objects `{}`

**Approach:** Implemented a tokenizer + state machine architecture with clear separation of concerns:
- **Tokenizer**: Breaks input into tokens (`LEFT_BRACE`, `RIGHT_BRACE`, `EOF`, `INVALID`)
- **State Machine**: Validates token sequence through defined states (`START` → `IN_OBJECT` → `COMPLETE`)

**Motivation:** Created an extensible foundation with proper tokenization that could be easily enhanced for future JSON constructs. The modular design separates lexical analysis from syntactic validation.

## Step 2: String Key-Value Pairs

**Objective:** Parse JSON objects with string keys and values like `{"key": "value"}` and `{"key1": "value1", "key2": "value2"}`

**Approach:** Replaced state machine with recursive descent parser while preserving the tokenizer foundation:
- **Enhanced Tokenizer**: Added `STRING`, `COLON`, `COMMA` tokens with proper string parsing and escape sequence handling
- **Recursive Descent Parser**: Grammar-driven approach with dedicated functions for each JSON construct (`parseObject`, `parseKeyValuePair`, `parseValue`)

**Motivation:** The recursive descent approach naturally mirrors JSON's grammar structure, making it much easier to extend for nested objects, arrays, and other complex structures in future steps. Each parser function has a single, clear responsibility aligned with JSON syntax rules.

## Step 3: Boolean, Null, and Numeric Values

**Objective:** Parse JSON objects with boolean, null, and numeric values like `{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`

**Approach:** Extended tokenizer and parser to handle primitive value types while maintaining the recursive descent architecture:
- **Enhanced Tokenizer**: Added `TRUE`, `FALSE`, `NULL`, `NUMBER` tokens with dedicated parsing methods (`parseKeywordToken`, `parseNumberToken`)
- **Character Dispatch**: Extended `NextToken()` to route letters (`t`, `f`, `n`) and digits (`0-9`) to appropriate parsing functions
- **Value Parser Extension**: Updated `parseValue()` to accept multiple token types with consistent error handling and position tracking

**Motivation:** Completing support for all JSON primitive value types creates a solid foundation before tackling complex nested structures. The token-based approach makes adding new value types straightforward while preserving the existing string and object parsing logic. Case-sensitive keyword matching ensures strict JSON compliance.

## Step 4: Arrays and Nested Objects

**Objective:** Parse JSON objects containing arrays and nested objects like `{"key-o": {"inner key": "inner value"}, "key-l": ["list value"]}`

**Approach:** Extended the parser with array support while leveraging the recursive descent architecture to enable unlimited nesting depth:
- **Enhanced Tokenizer**: Added `LEFT_BRACKET`, `RIGHT_BRACKET` tokens with simple character dispatch for `[` and `]`
- **Array Parser**: Implemented `parseArray()` function following the same comma-separated pattern as `parseObject()`
- **Recursive Value Parsing**: Enhanced `parseValue()` to handle objects and arrays as values, creating natural recursion
- **Unified Error Handling**: Extended trailing comma detection and position tracking to arrays

**Key Architecture Insight:** The recursive design proved its power - `parseValue()` can call `parseObject()` or `parseArray()`, which in turn call `parseValue()` for their contents, enabling structures like `{"a": [{"b": {"c": ["d"]}}]}` with no additional complexity.

**Motivation:** This step validates the architectural foundation by demonstrating that the recursive descent approach scales naturally to handle arbitrarily complex JSON structures. The clean separation between tokenization and parsing makes extending to new structural elements straightforward while maintaining all existing functionality.

## Usage

```bash
# Build the parser
go build -o json-parser

# Test with JSON files
./json-parser tests/step1/valid.json    # Empty object: {}
./json-parser tests/step2/valid.json    # Simple key-value pairs: {"key": "value"}
./json-parser tests/step3/valid.json    # Boolean, null, and numeric values
./json-parser tests/step4/valid.json    # Empty arrays and nested objects
./json-parser tests/step4/valid2.json   # Populated arrays and nested objects

# Test complex nested structures
echo '{"users": [{"name": "Alice", "scores": [95, 87]}, {"name": "Bob", "active": true}]}' | ./json-parser /dev/stdin
```

## Architecture Benefits

- **Clean Separation**: Tokenization and parsing phases are distinct and maintainable
- **Extensible Design**: Easy to add new token types and grammar rules (proven through 4 incremental steps)
- **Professional Foundation**: Industry-standard recursive descent approach scales to complex nested structures
- **Precise Error Reporting**: Detailed position and context information for all error cases
- **Natural Recursion**: Parser functions call each other recursively, enabling unlimited JSON nesting depth
- **Proven Scalability**: Successfully handles everything from `{}` to deeply nested objects and arrays
- **Robust Validation**: Comprehensive error detection including trailing commas, invalid syntax, and type mismatches

## Current Capabilities

✅ **All JSON Value Types**: strings, numbers, booleans, null, objects, arrays
✅ **Unlimited Nesting**: `{"a": [{"b": {"c": ["d"]}}]}` and beyond
✅ **Comprehensive Validation**: Trailing commas, quote types, case sensitivity
✅ **Professional Error Reporting**: Precise position information for debugging
✅ **Escape Sequences**: Full support for `\"`, `\\`, `\n`, `\t`, etc.
✅ **Whitespace Handling**: Proper normalization and formatting tolerance

## What's Next?

The parser foundation is solid and ready for advanced features:
- **Floating-point numbers**: `{"pi": 3.14159, "scientific": 1.23e-4}`
- **Negative numbers**: `{"temperature": -20, "balance": -1.50}`
- **Unicode strings**: `{"message": "Hello \u4e16\u754c"}`
- **Performance optimization**: Streaming, memory efficiency
- **AST generation**: Build data structures instead of just validation

## Project Structure

```
json-parser/
├── main.go                 # CLI entry point
├── parser/parser.go        # Complete tokenizer + parser implementation
├── tests/step[1-4]/        # Comprehensive test suite
├── go.mod                  # Go module definition
├── README.md              # This file
└── CLAUDE.md              # Technical implementation notes
```