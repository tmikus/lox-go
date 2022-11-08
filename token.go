package main

import "fmt"

type Token struct {
	Lexeme  string
	Line    uint
	Literal interface{}
	Type    TokenType
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line uint) Token {
	return Token{lexeme, line, literal, tokenType}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", t.Type, t.Lexeme, t.Literal)
}
