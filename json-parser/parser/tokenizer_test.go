package parser

import (
	"strings"
	"testing"
)

// Test NewTokenizer creation
func TestNewTokenizer(t *testing.T) {
	tokenizer := NewTokenizer(`{"test": "value"}`)
	if tokenizer == nil {
		t.Fatal("NewTokenizer returned nil")
	}
	if tokenizer.position != 0 {
		t.Errorf("Expected position 0, got %d", tokenizer.position)
	}
}

// Test NextToken with comprehensive token sequences
func TestNextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "empty object",
			input:    "{}",
			expected: []TokenType{LEFT_BRACE, RIGHT_BRACE, EOF},
		},
		{
			name:     "simple key-value",
			input:    `{"key": "value"}`,
			expected: []TokenType{LEFT_BRACE, STRING, COLON, STRING, RIGHT_BRACE, EOF},
		},
		{
			name:     "array with values",
			input:    `["val1", "val2"]`,
			expected: []TokenType{LEFT_BRACKET, STRING, COMMA, STRING, RIGHT_BRACKET, EOF},
		},
		{
			name:     "mixed values",
			input:    `{"bool": true, "num": 123, "null": null, "false": false}`,
			expected: []TokenType{LEFT_BRACE, STRING, COLON, TRUE, COMMA, STRING, COLON, NUMBER, COMMA, STRING, COLON, NULL, COMMA, STRING, COLON, FALSE, RIGHT_BRACE, EOF},
		},
		{
			name:     "nested structure",
			input:    `{"obj": {}, "arr": []}`,
			expected: []TokenType{LEFT_BRACE, STRING, COLON, LEFT_BRACE, RIGHT_BRACE, COMMA, STRING, COLON, LEFT_BRACKET, RIGHT_BRACKET, RIGHT_BRACE, EOF},
		},
		{
			name:     "whitespace handling",
			input:    "  {  \"key\"  :  \"value\"  }  ",
			expected: []TokenType{LEFT_BRACE, STRING, COLON, STRING, RIGHT_BRACE, EOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input)
			var tokens []TokenType

			for {
				token := tokenizer.NextToken()
				tokens = append(tokens, token.Type)
				if token.Type == EOF || token.Type == INVALID {
					break
				}
			}

			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", tokens)
				return
			}

			for i, expected := range tt.expected {
				if tokens[i] != expected {
					t.Errorf("Token %d: expected %s, got %s", i, expected, tokens[i])
				}
			}
		})
	}
}

// Test NextChar method
func TestNextChar(t *testing.T) {
	tokenizer := NewTokenizer("abc")

	if char := tokenizer.NextChar(); char != 'a' {
		t.Errorf("Expected 'a', got %c", char)
	}
	if tokenizer.position != 1 {
		t.Errorf("Expected position 1, got %d", tokenizer.position)
	}

	if char := tokenizer.NextChar(); char != 'b' {
		t.Errorf("Expected 'b', got %c", char)
	}

	if char := tokenizer.NextChar(); char != 'c' {
		t.Errorf("Expected 'c', got %c", char)
	}

	// Test EOF
	if char := tokenizer.NextChar(); char != 0 {
		t.Errorf("Expected EOF (0), got %c", char)
	}
}

// Test string parsing with all escape sequences
func TestParseStringToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		tokenType TokenType
	}{
		{"simple string", `"hello"`, "hello", STRING},
		{"empty string", `""`, "", STRING},
		{"string with spaces", `"hello world"`, "hello world", STRING},
		{"quote escape", `"say \"hello\""`, `say "hello"`, STRING},
		{"backslash escape", `"path\\to\\file"`, `path\to\file`, STRING},
		{"forward slash escape", `"url\/path"`, "url/path", STRING},
		{"backspace escape", `"text\btext"`, "text\btext", STRING},
		{"form feed escape", `"text\ftext"`, "text\ftext", STRING},
		{"newline escape", `"line1\nline2"`, "line1\nline2", STRING},
		{"carriage return escape", `"line1\rline2"`, "line1\rline2", STRING},
		{"tab escape", `"col1\tcol2"`, "col1\tcol2", STRING},
		{"all escapes", `"quote:\" slash:\\ forward:\/ back:\b form:\f new:\n ret:\r tab:\t"`,
			"quote:\" slash:\\ forward:/ back:\b form:\f new:\n ret:\r tab:\t", STRING},
		{"unterminated string", `"hello`, "unterminated string", INVALID},
		{"unterminated escape", `"hello\`, "unterminated string", INVALID},
		{"invalid escape", `"hello\x"`, "invalid escape sequence '\\x'", INVALID}, // Parser now rejects invalid escapes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenizer := NewTestingTokenizer(tt.input)
			// Skip the opening quote
			testTokenizer.SetPosition(1)
			token := testTokenizer.ParseStringToken(0)

			if token.Type != tt.tokenType {
				t.Errorf("Expected token type %s, got %s", tt.tokenType, token.Type)
			}
			if token.Value != tt.expected {
				t.Errorf("Expected value %q, got %q", tt.expected, token.Value)
			}
		})
	}
}

// Test keyword parsing with case sensitivity
func TestParseKeywordToken(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		firstChar rune
		expected  TokenType
		value     string
	}{
		{"true keyword", "true", 't', TRUE, "true"},
		{"false keyword", "false", 'f', FALSE, "false"},
		{"null keyword", "null", 'n', NULL, "null"},
		{"wrong case True", "True", 'T', INVALID, "True"},
		{"wrong case FALSE", "FALSE", 'F', INVALID, "FALSE"},
		{"wrong case NULL", "NULL", 'N', INVALID, "NULL"},
		{"partial true", "tr", 't', INVALID, "tr"},
		{"invalid keyword", "test", 't', INVALID, "test"},
		{"number start", "123", '1', INVALID, ""}, // Won't be called with digit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenizer := NewTestingTokenizer(tt.input)
			// Skip the first character since parseKeywordToken receives it
			testTokenizer.SetPosition(1)
			token := testTokenizer.ParseKeywordToken(0, tt.firstChar)

			if token.Type != tt.expected {
				t.Errorf("Expected token type %s, got %s", tt.expected, token.Type)
			}
			if tt.value != "" && token.Value != tt.value {
				t.Errorf("Expected value %q, got %q", tt.value, token.Value)
			}
		})
	}
}

// Test number parsing
func TestParseNumberToken(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		firstChar rune
		expected  string
	}{
		{"single digit", "1", '1', "1"},
		{"multiple digits", "123", '1', "123"},
		{"zero", "0", '0', "0"},
		{"large number", "999999", '9', "999999"},
		{"number followed by space", "42 ", '4', "42"},
		{"number followed by comma", "123,", '1', "123"},
		{"number followed by brace", "456}", '4', "456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTokenizer := NewTestingTokenizer(tt.input)
			// Skip the first character since parseNumberToken receives it
			testTokenizer.SetPosition(1)
			token := testTokenizer.ParseNumberToken(0, tt.firstChar)

			if token.Type != NUMBER {
				t.Errorf("Expected NUMBER token, got %s", token.Type)
			}
			if token.Value != tt.expected {
				t.Errorf("Expected value %q, got %q", tt.expected, token.Value)
			}
		})
	}
}

// Test whitespace handling
func TestWhitespaceHandling(t *testing.T) {
	inputs := []string{
		"   {   }   ",
		"\n{\n}\n",
		"\t{\t}\t",
		"\r{\r}\r",
		" \t\n\r { \t\n\r } \t\n\r ",
	}

	expected := []TokenType{LEFT_BRACE, RIGHT_BRACE, EOF}

	for _, input := range inputs {
		t.Run("whitespace: "+strings.ReplaceAll(input, "\n", "\\n"), func(t *testing.T) {
			tokenizer := NewTokenizer(input)
			var tokens []TokenType

			for {
				token := tokenizer.NextToken()
				tokens = append(tokens, token.Type)
				if token.Type == EOF || token.Type == INVALID {
					break
				}
			}

			if len(tokens) != len(expected) {
				t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
				return
			}

			for i, exp := range expected {
				if tokens[i] != exp {
					t.Errorf("Token %d: expected %s, got %s", i, exp, tokens[i])
				}
			}
		})
	}
}

// Test position tracking
func TestPositionTracking(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int // positions of tokens
	}{
		{"simple", `{}`, []int{0, 1, 2}},
		{"with spaces", ` { } `, []int{1, 3, 5}},
		{"string", `"test"`, []int{0, 6}},
		{"key-value", `"key":"value"`, []int{0, 5, 6, 13}},
		{"with whitespace", ` "key" : "value" `, []int{1, 7, 9, 17}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input)
			var positions []int

			for {
				token := tokenizer.NextToken()
				positions = append(positions, token.Position)
				if token.Type == EOF || token.Type == INVALID {
					break
				}
			}

			if len(positions) != len(tt.expected) {
				t.Errorf("Expected %d positions, got %d", len(tt.expected), len(positions))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", positions)
				return
			}

			for i, expected := range tt.expected {
				if positions[i] != expected {
					t.Errorf("Position %d: expected %d, got %d", i, expected, positions[i])
				}
			}
		})
	}
}

// Test invalid token detection
func TestInvalidTokenDetection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		position int
		value    string
	}{
		{"invalid character", "@", 0, "@"},
		{"invalid character mid", `{"key": @}`, 8, "@"},
		{"hash symbol", "#", 0, "#"},
		{"dollar sign", "$", 0, "$"},
		{"ampersand", "&", 0, "&"},
		{"question mark", "?", 0, "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input)
			var token Token

			// Get tokens until we find the invalid one
			for {
				token = tokenizer.NextToken()
				if token.Type == INVALID || token.Type == EOF {
					break
				}
			}

			if token.Type != INVALID {
				t.Errorf("Expected INVALID token, got %s", token.Type)
				return
			}

			if token.Position != tt.position {
				t.Errorf("Expected position %d, got %d", tt.position, token.Position)
			}

			if token.Value != tt.value {
				t.Errorf("Expected value %q, got %q", tt.value, token.Value)
			}
		})
	}
}

// Test TokenType String method
func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{LEFT_BRACE, "LEFT_BRACE"},
		{RIGHT_BRACE, "RIGHT_BRACE"},
		{LEFT_BRACKET, "LEFT_BRACKET"},
		{RIGHT_BRACKET, "RIGHT_BRACKET"},
		{STRING, "STRING"},
		{COLON, "COLON"},
		{COMMA, "COMMA"},
		{TRUE, "TRUE"},
		{FALSE, "FALSE"},
		{NULL, "NULL"},
		{NUMBER, "NUMBER"},
		{EOF, "EOF"},
		{INVALID, "INVALID"},
		{TokenType(999), "UNKNOWN"}, // Test unknown token type
	}

	for _, tt := range tests {
		if result := tt.tokenType.String(); result != tt.expected {
			t.Errorf("TokenType %d: expected %q, got %q", int(tt.tokenType), tt.expected, result)
		}
	}
}

// Test TestingTokenizer methods
func TestTestingTokenizerMethods(t *testing.T) {
	testTokenizer := NewTestingTokenizer("test input")

	// Test position getter/setter
	if pos := testTokenizer.GetPosition(); pos != 0 {
		t.Errorf("Expected initial position 0, got %d", pos)
	}

	testTokenizer.SetPosition(5)
	if pos := testTokenizer.GetPosition(); pos != 5 {
		t.Errorf("Expected position 5 after SetPosition, got %d", pos)
	}

	// Test exposed methods exist and can be called
	token := testTokenizer.ParseStringToken(0)
	if token.Type != INVALID { // Should be invalid since we're not starting at a quote
		t.Errorf("Expected INVALID token for non-quoted string")
	}
}

// Benchmark NextToken performance
func BenchmarkNextToken(b *testing.B) {
	input := `{"key1": "value1", "key2": 123, "key3": true, "key4": null, "key5": [1, 2, 3], "key6": {"nested": "object"}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer := NewTokenizer(input)
		for {
			token := tokenizer.NextToken()
			if token.Type == EOF || token.Type == INVALID {
				break
			}
		}
	}
}

// Benchmark string parsing
func BenchmarkParseStringToken(b *testing.B) {
	input := `"This is a test string with some \\"escaped\\" content and \nnewlines"`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testTokenizer := NewTestingTokenizer(input)
		testTokenizer.SetPosition(1) // Skip opening quote
		_ = testTokenizer.ParseStringToken(0)
	}
}