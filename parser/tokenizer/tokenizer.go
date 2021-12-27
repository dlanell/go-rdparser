package tokenizer

import (
	"errors"
	"fmt"
	"regexp"
)

type Tokenizer struct {
	text   string
	cursor int
}

type Props struct {
	Text string
}

type Token struct {
	TokenType string
	Value     string
}

const (
	NumberToken            string = "NUMBER"
	StringToken                   = "STRING"
	SemiColonToken                = ";"
	AdditiveOperator              = "+"
	MultiplicativeOperator        = "*"
	OpenCurlyBrace                = "{"
	CloseCurlyBrace               = "}"
	OpenParentheses               = "("
	CloseParentheses              = ")"
	Comma                         = ","
	RelationalOperator            = "RELATIONAL_OPERATOR"
	Identifier                    = "IDENTIFIER"
	SimpleAssignment              = "SIMPLE_ASSIGNMENT"
	ComplexAssignment             = "COMPLEX_ASSIGNMENT"
	LetKeyword                    = "let"
	IfKeyword                     = "if"
	ElseKeyword                   = "else"
	SkipToken                     = ""
)

func New(props Props) *Tokenizer {
	tokenizer := &Tokenizer{
		text:   props.Text,
		cursor: 0,
	}
	return tokenizer
}

func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.text)
}

func (t *Tokenizer) isEOF() bool {
	return t.cursor == len(t.text)
}

var spec = [][]string{
	//---------------------------------------------------
	// Whitespace

	{`^\s+`, SkipToken},

	//---------------------------------------------------
	// Comments

	// skip single-line comment
	{`^//.*`, SkipToken},
	// skip multi-line comment
	{`^/\*[\s\S]*?\*/`, SkipToken},

	//---------------------------------------------------
	// Symbols, Delimiters

	{`^;`, SemiColonToken},
	{`^{`, OpenCurlyBrace},
	{`^}`, CloseCurlyBrace},
	{`^\(`, OpenParentheses},
	{`^\)`, CloseParentheses},
	{`^\,`, Comma},

	//---------------------------------------------------
	// Keywords

	{`^\blet\b`, LetKeyword},
	{`^\bif\b`, IfKeyword},
	{`^\belse\b`, ElseKeyword},

	//---------------------------------------------------
	// Numbers

	{`^\d+`, NumberToken},

	//---------------------------------------------------
	// Identifiers

	{`^\w+`, Identifier},

	//---------------------------------------------------
	// Assignment operators =, +=, -=, *=, /=

	{`^=`, SimpleAssignment},
	{`^[+\-*/]=`, ComplexAssignment},

	//---------------------------------------------------
	// Math operators +, -, *, /

	{`^[+|-]`, AdditiveOperator},
	{`^[*|/]`, MultiplicativeOperator},

	//---------------------------------------------------
	// Relational operators >, >=, <, <=

	{`^[>|<]=?`, RelationalOperator},

	//---------------------------------------------------
	// Strings

	{`^"[^"]*"`, StringToken},
	{`^'[^']*'`, StringToken},
}

func (t *Tokenizer) GetNextToken() (*Token, error) {
	if !t.hasMoreTokens() {
		return nil, errors.New("no tokens present")
	}

	characters := []byte(t.text)[t.cursor:]

	for _, spec := range spec {
		regexText := spec[0]
		tokenType := spec[1]
		regex := regexp.MustCompile(regexText)
		tokenValue := t.match(regex, string(characters))
		if tokenValue == "" {
			continue
		}
		if tokenType == SkipToken {
			return t.GetNextToken()
		}
		return &Token{TokenType: tokenType, Value: tokenValue}, nil
	}

	return nil, fmt.Errorf(`unexpected token: %s`, string(characters[0]))
}

func (t *Tokenizer) match(regex *regexp.Regexp, text string) string {
	matchedToken := regex.FindString(text)
	if matchedToken == "" {
		return matchedToken
	}
	t.cursor += len(matchedToken)
	return matchedToken
}
