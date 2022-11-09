package main

import (
	"fmt"
	"strconv"
)

type ScannerError struct {
	Line    uint
	Message string
	Where   string
}

func (se ScannerError) Error() string {
	return fmt.Sprintf("[line %v] Error %v: %v", se.Line, se.Where, se.Message)
}

type Scanner struct {
	current uint
	line    uint
	source  string
	start   uint
}

func NewScanner(source string) Scanner {
	return Scanner{
		current: 0,
		line:    1,
		source:  source,
		start:   0,
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) createToken(tokenType TokenType) Token {
	// TODO: Implement 'literal'
	return NewToken(tokenType, s.source[s.start:s.current], nil, s.line)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= uint(len(s.source))
}

func (s *Scanner) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) parseNumber() (Option[Token], error) {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	// Look for a fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value := s.source[s.start:s.current]
	numericValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return NewEmptyOption[Token](), ScannerError{
			Line:    s.line,
			Message: "Invalid number.",
			Where:   "", // TODO: Implement
		}
	}
	return NewOption[Token](NewToken(NUMBER, s.source[s.start:s.current], numericValue, s.line)), nil
}

func (s *Scanner) parseString() (Option[Token], error) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return NewEmptyOption[Token](), ScannerError{
			Line:    s.line,
			Message: "Unterminated string.",
			Where:   "", // TODO: Implement
		}
	}
	// The closing "
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	return NewOption[Token](NewToken(STRING, value, value, s.line)), nil
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= uint(len(s.source)) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanToken() (Option[Token], error) {
	char := s.advance()
	switch char {
	case '(':
		return NewOption[Token](s.createToken(LEFT_PAREN)), nil
	case ')':
		return NewOption[Token](s.createToken(RIGHT_PAREN)), nil
	case '{':
		return NewOption[Token](s.createToken(LEFT_BRACE)), nil
	case '}':
		return NewOption[Token](s.createToken(RIGHT_BRACE)), nil
	case ',':
		return NewOption[Token](s.createToken(COMMA)), nil
	case '.':
		return NewOption[Token](s.createToken(DOT)), nil
	case '-':
		return NewOption[Token](s.createToken(MINUS)), nil
	case '+':
		return NewOption[Token](s.createToken(PLUS)), nil
	case ';':
		return NewOption[Token](s.createToken(SEMICOLON)), nil
	case '*':
		return NewOption[Token](s.createToken(STAR)), nil
	case '!':
		if s.match('=') {
			return NewOption[Token](s.createToken(BANG_EQUAL)), nil
		} else {
			return NewOption[Token](s.createToken(BANG)), nil
		}
	case '=':
		if s.match('=') {
			return NewOption[Token](s.createToken(EQUAL_EQUAL)), nil
		} else {
			return NewOption[Token](s.createToken(EQUAL)), nil
		}
	case '<':
		if s.match('=') {
			return NewOption[Token](s.createToken(LESS_EQUAL)), nil
		} else {
			return NewOption[Token](s.createToken(LESS)), nil
		}
	case '>':
		if s.match('=') {
			return NewOption[Token](s.createToken(GREATER_EQUAL)), nil
		} else {
			return NewOption[Token](s.createToken(GREATER)), nil
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
			return NewEmptyOption[Token](), nil
		} else {
			return NewOption[Token](s.createToken(SLASH)), nil
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
		return NewEmptyOption[Token](), nil
	case '\n':
		s.line++
		return NewEmptyOption[Token](), nil
	case '"':
		return s.parseString()
	}
	if s.isDigit(char) {
		return s.parseNumber()
	}
	return NewEmptyOption[Token](), ScannerError{
		Line:    s.line,
		Message: "Unexpected character: " + string(char),
		Where:   "", // TODO: Implement
	}
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	var tokens []Token
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		token, err := s.scanToken()
		if err != nil {
			return nil, err
		}
		if token.IsEmpty() {
			continue
		}
		tokens = append(tokens, token.Value())
	}
	tokens = append(tokens, NewToken(EOF, "", nil, s.line))
	return tokens, nil
}
