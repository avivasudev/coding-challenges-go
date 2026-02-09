# JSON Parser

## Overview

A step-by-step implementation of a JSON parser in Go, built incrementally to support increasingly complex JSON structures. The project emphasizes clean architecture, extensibility, and proper error handling.

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

## Usage

```bash
# Build the parser
go build -o json-parser

# Test with JSON files
./json-parser tests/step1/valid.json    # Empty object
./json-parser tests/step2/valid.json    # Simple key-value pairs
```

## Architecture Benefits

- **Clean Separation**: Tokenization and parsing phases are distinct
- **Extensible Design**: Easy to add new token types and grammar rules
- **Professional Foundation**: Industry-standard recursive descent approach
- **Precise Error Reporting**: Detailed position and context information