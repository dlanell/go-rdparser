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

type Spec map[*regexp.Regexp]string

var spec = Spec{
	//---------------------------------------------------
	// Whitespace

	regexp.MustCompile(`^\s+`): SkipToken,

	//---------------------------------------------------
	// Comments

	// skip single-line comment
	regexp.MustCompile(`^//.*`): SkipToken,
	// skip multi-line comment
	regexp.MustCompile(`^/\*[\s\S]*?\*/`): SkipToken,

	//---------------------------------------------------
	// Symbols, Delimiters

	regexp.MustCompile(`^;`):  SemiColonToken,
	regexp.MustCompile(`^{`):  OpenCurlyBrace,
	regexp.MustCompile(`^}`):  CloseCurlyBrace,
	regexp.MustCompile(`^\(`): OpenParentheses,
	regexp.MustCompile(`^\)`): CloseParentheses,

	//---------------------------------------------------
	// Math operators +, -, *, /

	regexp.MustCompile(`^[+|-]`): AdditiveOperator,
	regexp.MustCompile(`^[*|/]`): MultiplicativeOperator,

	//---------------------------------------------------
	// Numbers

	regexp.MustCompile(`^\d+`): NumberToken,

	//---------------------------------------------------
	// Strings

	regexp.MustCompile(`^"[^"]*"`): StringToken,
	regexp.MustCompile(`^'[^']*'`): StringToken,
}

func (t *Tokenizer) GetNextToken() (*Token, error) {
	if !t.hasMoreTokens() {
		return nil, errors.New("no tokens present")
	}

	characters := []byte(t.text)[t.cursor:]

	for regex, tokenType := range spec {
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
