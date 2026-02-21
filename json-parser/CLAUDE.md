# JSON Parser Project - Claude Reference

## Project Overview

A **step-by-step JSON parser implementation in Go** that incrementally adds support for increasingly complex JSON structures. The project emphasizes clean architecture, extensibility, and proper error handling.

## Current Status

### ✅ Completed Steps
- **Step 1**: Parse empty objects `{}`
- **Step 2**: Parse string key-value pairs `{"key": "value"}` and multiple pairs
- **Step 3**: Parse boolean, null, and numeric values `{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`
- **Step 4**: Parse arrays and nested objects `{"key-o": {"inner key": "inner value"}, "key-l": ["list value"]}`

### 🎯 Current Capabilities
- Empty objects: `{}`
- Single key-value pairs: `{"key": "value"}`
- Multiple key-value pairs: `{"key1": "value1", "key2": "value2"}`
- **Boolean values**: `true`, `false` (case-sensitive)
- **Null values**: `null`
- **Numeric values**: positive integers like `101`
- **Arrays**: empty `[]` and populated `["value1", "value2"]`
- **Nested objects**: `{"outer": {"inner": "value"}}`
- **Mixed structures**: objects containing arrays, arrays containing objects
- **Arbitrarily deep nesting**: `{"a": [{"b": {"c": ["d"]}}]}`
- Whitespace handling and normalization
- String escape sequence support (`\"`, `\\`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t`)
- Trailing comma detection and rejection (objects and arrays)
- Detailed error reporting with position information

## Architecture

### Two-Phase Design
1. **Tokenizer** (`parser/parser.go:51-200`):
   - Lexical analysis: breaks input into tokens
   - Supported tokens: `LEFT_BRACE`, `RIGHT_BRACE`, `LEFT_BRACKET`, `RIGHT_BRACKET`, `STRING`, `COLON`, `COMMA`, `TRUE`, `FALSE`, `NULL`, `NUMBER`, `EOF`, `INVALID`
   - Position tracking for error reporting
   - String parsing with escape sequence handling
   - Keyword parsing with case-sensitive matching
   - Integer number parsing

2. **Recursive Descent Parser** (`parser/parser.go:240-390`):
   - Syntactic analysis: validates token sequences against JSON grammar
   - Key functions:
     - `parseObject()` - handles object structure `{ ... }`
     - `parseArray()` - handles array structure `[ ... ]`
     - `parseKeyValuePair()` - handles `"key": "value"` pairs
     - `parseValue()` - supports all JSON value types (primitives, objects, arrays)
   - Grammar-driven approach that mirrors JSON structure
   - **Recursive by design**: `parseValue()` can call `parseObject()` or `parseArray()`, enabling unlimited nesting depth

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
│   └── step4/             # Array and nested object tests
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

# Step 4: Arrays and nested objects
./json-parser tests/step4/valid.json      # Should print "Valid JSON"
./json-parser tests/step4/valid2.json     # Should print "Valid JSON"
./json-parser tests/step4/invalid.json    # Should print error
```

### Test Files Content
- `tests/step1/valid.json`: `{}`
- `tests/step2/valid.json`: `{"key": "value"}`
- `tests/step2/valid2.json`: `{"key": "value", "key2": "value"}`
- `tests/step3/valid.json`: `{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`
- `tests/step3/invalid.json`: `{"key2": False}` (case sensitivity test - "False" should be rejected)
- `tests/step4/valid.json`: `{"key": "value", "key-n": 101, "key-o": {}, "key-l": []}` (empty objects and arrays)
- `tests/step4/valid2.json`: `{"key": "value", "key-n": 101, "key-o": {"inner key": "inner value"}, "key-l": ["list value"]}` (nested objects and arrays with values)
- `tests/step4/invalid.json`: `{"key-l": ['list value']}` (single quotes in array - should be rejected)

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

## ✅ Step 4 Implementation (Completed)

**Successfully implemented support for:**
1. **Arrays**: Empty `[]` and populated arrays `["value1", "value2"]`
2. **Nested objects**: Objects as values `{"outer": {"inner": "value"}}`
3. **Mixed structures**: Objects containing arrays, arrays containing objects
4. **Arbitrary nesting depth**: `{"a": [{"b": {"c": ["d"]}}]}`

**Implementation details:**
- Added new token types: `LEFT_BRACKET`, `RIGHT_BRACKET`
- Updated `NextToken()` character dispatch for `[` and `]`
- Implemented `parseArray()` function with comma-separated value parsing
- Enhanced `parseValue()` to handle objects and arrays as values
- Recursive design: `parseValue()` → `parseObject()` → `parseValue()` enables unlimited nesting
- Maintained trailing comma detection and rejection for arrays
- Preserved all error reporting and position tracking

**Key Architecture Decision:**
The recursive descent approach proved ideal for nested structures. By allowing `parseValue()` to call both `parseObject()` and `parseArray()`, the parser naturally handles arbitrarily complex nesting without additional complexity.

## Next Steps (Step 5)

Based on typical JSON parser evolution, Step 5 might include:

1. **Floating-point numbers**: `{"pi": 3.14159, "scientific": 1.23e-4}`
2. **Negative numbers**: `{"temperature": -20, "balance": -1.50}`
3. **Unicode strings**: `{"unicode": "Hello \u4e16\u754c"}`
4. **More robust number parsing**: Handle edge cases, overflow, precision

### Architecture Benefits for Extension
- Number parsing logic is already isolated in `parseNumberToken()`
- Unicode support only requires extending string parsing
- Token-based approach makes extending numeric formats straightforward
- Error reporting system ready for number format validation

## Implementation Patterns (Critical for Resumption)

### Adding New Token Types (Step 3 & 4 Pattern)
1. **Extend TokenType enum** (around line 11):
   ```go
   const (
       // ... existing types
       NEW_TOKEN_TYPE  // Add here, before EOF
   )
   ```
2. **Add String() method case** (around line 40):
   ```go
   case NEW_TOKEN_TYPE:
       return "NEW_TOKEN_TYPE"
   ```
3. **Add NextToken() dispatch** (around line 230):
   ```go
   case 'x', 'y':  // triggering characters
       return t.parseNewTokenType(tokenPos, char)
   ```
4. **Create parsing method** (after parseNumberToken):
   ```go
   func (t *Tokenizer) parseNewTokenType(startPos int, firstChar rune) Token
   ```
5. **Update parseValue()** to accept new token type

### Adding New Structural Elements (Step 4 Pattern)
For composite structures like arrays, the pattern is:
1. **Add token types** for delimiters (`LEFT_BRACKET`, `RIGHT_BRACKET`)
2. **Add simple character dispatch** (no complex parsing needed for delimiters)
3. **Create parsing function** that mirrors `parseObject()` pattern:
   ```go
   func (p *Parser) parseArray() error {
       // Expect opening delimiter
       // Handle empty case
       // Parse first element
       // Loop for comma-separated additional elements
       // Check for trailing comma (reject)
       // Expect closing delimiter
   }
   ```
4. **Extend parseValue()** to handle new structure type
5. **Leverage recursion**: new parser can call `parseValue()`, enabling nesting

### Key Function Signatures
- `func NewTokenizer(input string) *Tokenizer`
- `func (t *Tokenizer) NextToken() Token`
- `func (t *Tokenizer) parseKeywordToken(startPos int, firstChar rune) Token`
- `func (t *Tokenizer) parseNumberToken(startPos int, firstChar rune) Token`
- `func NewParser(input string) *Parser`
- `func (p *Parser) ParseJSON() error`
- `func (p *Parser) parseObject() error`
- `func (p *Parser) parseArray() error` (added in Step 4)
- `func (p *Parser) parseValue() error`
- `func ValidateJSON(input string) error` (main entry point)

### Error Message Pattern
Always use: `fmt.Errorf("message at position %d", p.currentToken.Position)`

### Testing Pattern
1. Add test files to appropriate `tests/stepX/` directory
2. Test both valid and invalid cases
3. Verify error messages include position information
4. Run regression tests: `./json-parser tests/step1/valid.json` etc.

## Notes for Future Development

- **✅ Proven Architecture**: The recursive descent approach successfully handles nested objects and arrays with unlimited depth
- Consider adding an AST (Abstract Syntax Tree) representation for step 5+ if data structure building is needed
- Current parser validates but doesn't build a data structure (validation-only) - this is by design and works well
- Position tracking enables precise error messages for debugging complex nested structures
- Modular design allows independent testing of tokenizer vs parser logic
- **Recursive Value Parsing**: The key insight from Step 4 - `parseValue()` calling `parseObject()` and `parseArray()` creates natural recursion
- **Critical**: Always read the existing parser.go file first to understand current line numbers and structure
- **Step 4 Success**: Arrays and nested objects work flawlessly, proving the architecture scales