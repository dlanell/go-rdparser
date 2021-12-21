package tokenizer

import "strconv"

type Tokenizer struct {
	text   string
	cursor int
}

type Props struct {
	text string
}

type Token struct {
	tokenType string
	value     string
}

const (
	NumberToken string = "Number"
)

func New(props Props) *Tokenizer {
	tokenizer := &Tokenizer{
		text:   props.text,
		cursor: 0,
	}
	return tokenizer
}

func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.text)
}

func (t *Tokenizer) GetNextToken() *Token {
	if !t.hasMoreTokens() {
		return nil
	}
	characters := []rune(t.text)[t.cursor:]

	if _, charErr := strconv.Atoi(string(characters[0])); charErr == nil {
		tokenText := ""

		for ok := true; ok; ok = charErr == nil && t.hasMoreTokens() {
			tokenText += string(characters[t.cursor])
			t.cursor += 1
			if !t.hasMoreTokens() {
				break
			}
			_, charErr = strconv.Atoi(string(characters[t.cursor]))
		}
		return &Token{tokenType: NumberToken, value: tokenText}
	}
	return nil
}
