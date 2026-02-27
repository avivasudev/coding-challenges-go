package parser

import (
	"fmt"
	"unicode"
)

// JSONError provides structured error information for better testing
type JSONError struct {
	Message   string
	Position  int
	TokenType TokenType
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("%s at position %d", e.Message, e.Position)
}

// TokenType represents different types of tokens
type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
	STRING
	COLON
	COMMA
	TRUE
	FALSE
	NULL
	NUMBER
	EOF
	INVALID
)

// Token represents a single token with its type, value, and position
type Token struct {
	Type     TokenType
	Value    string
	Position int
}

// String method for better debugging
func (t TokenType) String() string {
	switch t {
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case LEFT_BRACKET:
		return "LEFT_BRACKET"
	case RIGHT_BRACKET:
		return "RIGHT_BRACKET"
	case STRING:
		return "STRING"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case NUMBER:
		return "NUMBER"
	case EOF:
		return "EOF"
	case INVALID:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}

// Tokenizer breaks input string into tokens
type Tokenizer struct {
	input    string
	position int
}

// NewTokenizer creates a new tokenizer with the given input
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:    input,
		position: 0,
	}
}

// NextChar returns the current character and advances position
func (t *Tokenizer) NextChar() rune {
	if t.position >= len(t.input) {
		return 0 // EOF
	}
	char := rune(t.input[t.position])
	t.position++
	return char
}


// skipWhitespace skips all whitespace characters
func (t *Tokenizer) skipWhitespace() {
	for t.position < len(t.input) {
		char := rune(t.input[t.position])
		if unicode.IsSpace(char) {
			t.position++
		} else {
			break
		}
	}
}

// parseStringToken reads a complete string token with escape handling
func (t *Tokenizer) parseStringToken(startPos int) Token {
	var result string

	for t.position < len(t.input) {
		char := rune(t.input[t.position])

		if char == '"' {
			// End of string
			t.position++
			return Token{Type: STRING, Value: result, Position: startPos}
		}

		if char == '\\' {
			// Handle escape sequences
			t.position++
			if t.position >= len(t.input) {
				return Token{Type: INVALID, Value: "unterminated string", Position: startPos}
			}

			nextChar := rune(t.input[t.position])
			switch nextChar {
			case '"':
				result += "\""
			case '\\':
				result += "\\"
			case '/':
				result += "/"
			case 'b':
				result += "\b"
			case 'f':
				result += "\f"
			case 'n':
				result += "\n"
			case 'r':
				result += "\r"
			case 't':
				result += "\t"
			default:
				result += string(nextChar)
			}
			t.position++
		} else {
			result += string(char)
			t.position++
		}
	}

	// If we reach here, string was not terminated
	return Token{Type: INVALID, Value: "unterminated string", Position: startPos}
}

// parseKeywordToken reads a complete keyword token (true, false, null)
func (t *Tokenizer) parseKeywordToken(startPos int, firstChar rune) Token {
	var keyword string
	keyword += string(firstChar)

	// Read alphabetic characters
	for t.position < len(t.input) {
		char := rune(t.input[t.position])
		if unicode.IsLetter(char) {
			keyword += string(char)
			t.position++
		} else {
			break
		}
	}

	// Match against valid keywords (case-sensitive)
	switch keyword {
	case "true":
		return Token{Type: TRUE, Value: keyword, Position: startPos}
	case "false":
		return Token{Type: FALSE, Value: keyword, Position: startPos}
	case "null":
		return Token{Type: NULL, Value: keyword, Position: startPos}
	default:
		return Token{Type: INVALID, Value: keyword, Position: startPos}
	}
}

// parseNumberToken reads a complete number token (positive integers)
func (t *Tokenizer) parseNumberToken(startPos int, firstChar rune) Token {
	var number string
	number += string(firstChar)

	// Check for invalid leading zeros (JSON spec: numbers cannot have leading zeros except for "0")
	if firstChar == '0' {
		// If we start with '0', only allow single '0' or immediately stop
		if t.position < len(t.input) {
			nextChar := rune(t.input[t.position])
			if unicode.IsDigit(nextChar) {
				// Leading zero followed by another digit is invalid (like "01", "013")
				return Token{Type: INVALID, Value: "numbers cannot have leading zeros", Position: startPos}
			}
		}
		return Token{Type: NUMBER, Value: number, Position: startPos}
	}

	// Read consecutive digits for non-zero numbers
	for t.position < len(t.input) {
		char := rune(t.input[t.position])
		if unicode.IsDigit(char) {
			number += string(char)
			t.position++
		} else {
			break
		}
	}

	return Token{Type: NUMBER, Value: number, Position: startPos}
}

// NextToken returns the next token from the input
func (t *Tokenizer) NextToken() Token {
	// Skip any leading whitespace
	t.skipWhitespace()

	// Remember position for token
	tokenPos := t.position

	// Get current character
	char := t.NextChar()

	// Process the character
	switch char {
	case 0: // EOF
		return Token{Type: EOF, Value: "", Position: tokenPos}
	case '{':
		return Token{Type: LEFT_BRACE, Value: "{", Position: tokenPos}
	case '}':
		return Token{Type: RIGHT_BRACE, Value: "}", Position: tokenPos}
	case '[':
		return Token{Type: LEFT_BRACKET, Value: "[", Position: tokenPos}
	case ']':
		return Token{Type: RIGHT_BRACKET, Value: "]", Position: tokenPos}
	case '"':
		// Parse string token (don't include the quote)
		return t.parseStringToken(tokenPos)
	case ':':
		return Token{Type: COLON, Value: ":", Position: tokenPos}
	case ',':
		return Token{Type: COMMA, Value: ",", Position: tokenPos}
	case 't', 'f', 'n':
		return t.parseKeywordToken(tokenPos, char)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return t.parseNumberToken(tokenPos, char)
	default:
		// Any other character is invalid
		return Token{Type: INVALID, Value: string(char), Position: tokenPos}
	}
}

// Parser structure that wraps the tokenizer and tracks current token
type Parser struct {
	tokenizer    *Tokenizer
	currentToken Token
	position     int
}

// NewParser creates a new parser with the given input
func NewParser(input string) *Parser {
	tokenizer := NewTokenizer(input)
	parser := &Parser{
		tokenizer: tokenizer,
		position:  0,
	}
	parser.advance() // Load first token
	return parser
}

// advance moves to the next token
func (p *Parser) advance() {
	p.currentToken = p.tokenizer.NextToken()
	p.position++
}

// ParseJSON is the main entry point for parsing
func (p *Parser) ParseJSON() error {
	err := p.parseValue()  // Accept any JSON value, not just objects
	if err != nil {
		return err
	}

	if p.currentToken.Type != EOF {
		return fmt.Errorf("unexpected token after JSON at position %d", p.currentToken.Position)
	}

	return nil
}

// parseObject handles { key:value, key:value }
func (p *Parser) parseObject() error {
	if p.currentToken.Type != LEFT_BRACE {
		return fmt.Errorf("expected '{' at position %d", p.currentToken.Position)
	}
	p.advance()

	// Handle empty object
	if p.currentToken.Type == RIGHT_BRACE {
		p.advance()
		return nil
	}

	// Parse first key-value pair
	err := p.parseKeyValuePair()
	if err != nil {
		return err
	}

	// Parse additional key-value pairs
	for p.currentToken.Type == COMMA {
		p.advance()

		// Check for trailing comma (invalid)
		if p.currentToken.Type == RIGHT_BRACE {
			return fmt.Errorf("trailing comma is not allowed at position %d", p.currentToken.Position)
		}

		err := p.parseKeyValuePair()
		if err != nil {
			return err
		}
	}

	if p.currentToken.Type != RIGHT_BRACE {
		return fmt.Errorf("expected '}' at position %d", p.currentToken.Position)
	}
	p.advance()

	return nil
}

// parseKeyValuePair handles "key": "value"
func (p *Parser) parseKeyValuePair() error {
	// Parse key
	if p.currentToken.Type != STRING {
		return fmt.Errorf("expected string key at position %d", p.currentToken.Position)
	}
	p.advance()

	// Parse colon
	if p.currentToken.Type != COLON {
		return fmt.Errorf("expected ':' after key at position %d", p.currentToken.Position)
	}
	p.advance()

	// Parse value
	return p.parseValue()
}

// parseValue handles string, boolean, null, number, array, and object values (for Step 4)
func (p *Parser) parseValue() error {
	switch p.currentToken.Type {
	case STRING, TRUE, FALSE, NULL, NUMBER:
		p.advance()
		return nil
	case LEFT_BRACE:
		return p.parseObject()
	case LEFT_BRACKET:
		return p.parseArray()
	case INVALID:
		// Return the specific error message from the tokenizer
		return fmt.Errorf("%s at position %d", p.currentToken.Value, p.currentToken.Position)
	default:
		return fmt.Errorf("expected value at position %d", p.currentToken.Position)
	}
}

// parseArray handles [ value, value, value ]
func (p *Parser) parseArray() error {
	if p.currentToken.Type != LEFT_BRACKET {
		return fmt.Errorf("expected '[' at position %d", p.currentToken.Position)
	}
	p.advance()

	// Handle empty array
	if p.currentToken.Type == RIGHT_BRACKET {
		p.advance()
		return nil
	}

	// Parse first value
	err := p.parseValue()
	if err != nil {
		return err
	}

	// Parse additional values
	for p.currentToken.Type == COMMA {
		p.advance()

		// Check for trailing comma (invalid)
		if p.currentToken.Type == RIGHT_BRACKET {
			return fmt.Errorf("trailing comma is not allowed at position %d", p.currentToken.Position)
		}

		err := p.parseValue()
		if err != nil {
			return err
		}
	}

	if p.currentToken.Type != RIGHT_BRACKET {
		return fmt.Errorf("expected ']' at position %d", p.currentToken.Position)
	}
	p.advance()

	return nil
}

// ValidateJSON validates if the input string is valid JSON
func ValidateJSON(input string) error {
	parser := NewParser(input)
	return parser.ParseJSON()
}

// TestingTokenizer provides access to tokenizer internals for testing
type TestingTokenizer struct {
	*Tokenizer
}

// NewTestingTokenizer creates a tokenizer exposed for testing
func NewTestingTokenizer(input string) *TestingTokenizer {
	return &TestingTokenizer{NewTokenizer(input)}
}

// Expose private methods for unit testing
func (tt *TestingTokenizer) ParseStringToken(startPos int) Token {
	return tt.parseStringToken(startPos)
}

func (tt *TestingTokenizer) ParseKeywordToken(startPos int, firstChar rune) Token {
	return tt.parseKeywordToken(startPos, firstChar)
}

func (tt *TestingTokenizer) ParseNumberToken(startPos int, firstChar rune) Token {
	return tt.parseNumberToken(startPos, firstChar)
}

func (tt *TestingTokenizer) SetPosition(pos int) { tt.position = pos }
func (tt *TestingTokenizer) GetPosition() int { return tt.position }