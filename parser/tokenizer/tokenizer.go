package tokenizer

import (
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
	NumberToken string = "NUMBER"
	StringToken        = "STRING"
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

func (t *Tokenizer) GetNextToken() *Token {
	if !t.hasMoreTokens() {
		return nil
	}

	characters := []byte(t.text)[t.cursor:]

	numberRegex := regexp.MustCompile(`^\d+`)
	if matchedToken := numberRegex.FindString(string(characters)); matchedToken != "" {
		t.cursor += len(matchedToken)
		return &Token{TokenType: NumberToken, Value: matchedToken}
	}

	doubleQuoteStringRegex := regexp.MustCompile(`^"[^"]*"`)
	if matchedToken := doubleQuoteStringRegex.FindString(string(characters)); matchedToken != ""{
		t.cursor += len(matchedToken)
		return &Token{TokenType: StringToken, Value: matchedToken}
	}

	singleQuoteStringRegex := regexp.MustCompile(`^'[^']*'`)
	if matchedToken := singleQuoteStringRegex.FindString(string(characters)); matchedToken != ""{
		t.cursor += len(matchedToken)
		return &Token{TokenType: StringToken, Value: matchedToken}
	}

	return nil
}
