package parser

import (
	"fmt"
	"unicode"
)

// TokenType represents different types of tokens
type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	STRING
	COLON
	COMMA
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
	case STRING:
		return "STRING"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
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
	case '"':
		// Parse string token (don't include the quote)
		return t.parseStringToken(tokenPos)
	case ':':
		return Token{Type: COLON, Value: ":", Position: tokenPos}
	case ',':
		return Token{Type: COMMA, Value: ",", Position: tokenPos}
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
	err := p.parseObject()
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

// parseValue handles string values (for Step 2)
func (p *Parser) parseValue() error {
	if p.currentToken.Type != STRING {
		return fmt.Errorf("expected string value at position %d", p.currentToken.Position)
	}
	p.advance()
	return nil
}

// ValidateJSON validates if the input string is valid JSON
func ValidateJSON(input string) error {
	parser := NewParser(input)
	return parser.ParseJSON()
}