package main

import (
	"fmt"
	"os"
	"strconv"
)

type TokenType int

func (tt TokenType) String() string {
	switch tt {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

const (
	// single char tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	BANG

	// one or more char tokens
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// arbitrarily long tokens
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	// special
	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Lexeme: %s, Literal: %v, Line: %d}", t.Type, t.Lexeme, t.Literal, t.Line)
}

type Scanner struct {
	source string
	tokens []Token

	line    int
	start   int
	current int
}

func NewScanner(code string) *Scanner {
	return &Scanner{
		source:  code,
		tokens:  make([]Token, 0),
		line:    1,
		start:   0,
		current: 0,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case ' ', '\r', '\t':
		break
	case '{':
		s.addTokenSingle(LEFT_BRACE)
	case '}':
		s.addTokenSingle(RIGHT_PAREN)
	case '(':
		s.addTokenSingle(LEFT_BRACE)
	case ')':
		s.addTokenSingle(RIGHT_BRACE)
	case ',':
		s.addTokenSingle(COMMA)
	case '.':
		s.addTokenSingle(DOT)
	case '-':
		s.addTokenSingle(MINUS)
	case '+':
		s.addTokenSingle(PLUS)
	case ';':
		s.addTokenSingle(SEMICOLON)
	case '*':
		s.addTokenSingle(STAR)
	case '!':
		s.addTokenPeek('=', BANG_EQUAL, EQUAL)
	case '=':
		s.addTokenPeek('=', EQUAL_EQUAL, EQUAL)
	case '<':
		s.addTokenPeek('=', LESS_EQUAL, LESS)
	case '>':
		s.addTokenPeek('=', GREATER_EQUAL, GREATER)
	case '/':
		switch s.peek() {
		case '/':
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		case '*':
			s.consumeMultilineComment()
		default:
			s.addTokenSingle(SLASH)
		}

	case '"':
		s.consumeString()
	case '\n':
		s.line++
	default:
		if IsDigit(c) {
			s.consumeNumber()
		} else if IsAlpha(c) {
			s.consumeIndentifier()
		} else {
			s.reportError(fmt.Sprintf("Unexpected character: %c", c))
			break
		}
	}
}

func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func IsAlphaNumeric(c rune) bool {
	return IsDigit(c) || IsAlpha(c)
}

func (s *Scanner) consumeNumber() {
	for IsDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && IsDigit(s.peekNext()) {
		s.advance()
		for IsDigit(s.peek()) {
			s.advance()
		}
	}
	value, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(NUMBER, value)
}

func (s *Scanner) consumeString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.reportError("Error unterminated string")
		return
	}
	s.advance() // The closing ".
	// Trim the surrounding quotes.
	s.addToken(STRING, s.source[s.start+1:s.current-1])
}

func (s *Scanner) consumeMultilineComment() {
	for !s.isAtEnd() {
		if s.peek() == '*' && s.peekNext() == '/' {
			s.advance()
			s.advance()
			break
		}
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.reportError("Unterminated multiline comment")
	}
}

func (s *Scanner) consumeIndentifier() {
	for IsAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	if tt, ok := keywords[text]; ok {
		s.addToken(tt, text)
	} else {
		s.addToken(IDENTIFIER, text)
	}
}

func (s *Scanner) reportError(msg string) {
	fmt.Printf("Error on line %d: %s\n", s.line, msg)
	os.Exit(1)
}

func (s *Scanner) addTokenSingle(tt TokenType) {
	s.addToken(tt, nil)
}

func (s *Scanner) addTokenPeek(expect rune, dual TokenType, single TokenType) {
	if s.peek() == expect {
		s.current++
		s.addTokenSingle(dual)
	} else {
		s.addTokenSingle(single)
	}
}

func (s *Scanner) addToken(tt TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{tt, s.source[s.start:s.current], literal, s.line})
}

func (s *Scanner) advance() rune {
	pop := rune(s.source[s.current])
	s.current++
	return pop
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return rune(s.source[s.current+1])
}
