package querytokenizer

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
	NumberToken        string = "NUMBER"
	StringToken               = "STRING"
	BooleanToken              = "BOOLEAN"
	OpenParentheses           = "("
	CloseParentheses          = ")"
	Comma                     = ","
	RelationalOperator        = "RELATIONAL_OPERATOR"
	LogicalOperator           = "LOGICAL_OPERATOR"
	Identifier                = "IDENTIFIER"
	SkipToken                 = ""
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
	// Symbols, Delimiters

	{`^\(`, OpenParentheses},
	{`^\)`, CloseParentheses},
	{`^\,`, Comma},

	//---------------------------------------------------
	// Numbers

	{`^\d+`, NumberToken},

	//---------------------------------------------------
	// Relational operators
	// eq -> equal
	// ne -> not equal
	// gt -> greater than
	// ge -> greater than or equal
	// lt -> less than
	// le -> less than or equal

	{`^(\beq\b|\bne\b|\blt\b|\ble\b|\bgt\b|\bge\b)`, RelationalOperator},

	//---------------------------------------------------
	// logical operators and, or

	{`^(\band\b|\bor\b)`, LogicalOperator},

	//---------------------------------------------------
	// Boolean value: true, false

	{`^(\btrue\b|\bfalse\b)`, BooleanToken},

	//---------------------------------------------------
	// Identifiers

	{`^\w+`, Identifier},

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
