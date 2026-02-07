package parser

import (
	"errors"
	"unicode"
)

// TokenType represents different types of tokens
type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
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
	default:
		// Any other character is invalid for step 1
		return Token{Type: INVALID, Value: string(char), Position: tokenPos}
	}
}

// Parser state constants
type ParserState int

const (
	START ParserState = iota
	IN_OBJECT
	COMPLETE
	ERROR
)

// String method for better debugging
func (s ParserState) String() string {
	switch s {
	case START:
		return "START"
	case IN_OBJECT:
		return "IN_OBJECT"
	case COMPLETE:
		return "COMPLETE"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ValidateJSON validates if the input string is valid JSON for step 1
func ValidateJSON(input string) error {
	tokenizer := NewTokenizer(input)
	state := START

	for {
		token := tokenizer.NextToken()

		// Process token based on current state
		switch state {
		case START:
			switch token.Type {
			case LEFT_BRACE:
				state = IN_OBJECT
			case EOF:
				return errors.New("empty input - not valid JSON")
			default:
				return errors.New("expected '{' at start of JSON")
			}

		case IN_OBJECT:
			switch token.Type {
			case RIGHT_BRACE:
				state = COMPLETE
			default:
				return errors.New("expected '}' after '{'")
			}

		case COMPLETE:
			switch token.Type {
			case EOF:
				return nil // Success!
			default:
				return errors.New("unexpected characters after complete JSON")
			}
		}
	}
}