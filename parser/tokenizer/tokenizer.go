package tokenizer

import "strconv"

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
	characters := []rune(t.text)[t.cursor:]

	if _, charErr := strconv.Atoi(string(characters[0])); charErr == nil {
		tokenText := ""

		for ok := true; ok; ok = charErr == nil && t.hasMoreTokens() {
			tokenText += string(characters[t.cursor])
			t.cursor++
			if !t.hasMoreTokens() {
				break
			}
			_, charErr = strconv.Atoi(string(characters[t.cursor]))
		}
		return &Token{TokenType: NumberToken, Value: tokenText}
	}

	if string(characters[0]) == `"` {
		tokenText := ``

		for ok := true; ok; ok = string(characters[t.cursor]) != `"` && !t.isEOF() {
			tokenText += string(characters[t.cursor])
			t.cursor++
		}
		tokenText += string(characters[t.cursor])
		t.cursor++

		return &Token{TokenType: StringToken, Value: tokenText}
	}

	if string(characters[0]) == `'` {
		tokenText := ``

		for ok := true; ok; ok = string(characters[t.cursor]) != `'` && !t.isEOF() {
			tokenText += string(characters[t.cursor])
			t.cursor++
		}
		tokenText += string(characters[t.cursor])
		t.cursor++

		return &Token{TokenType: StringToken, Value: tokenText}
	}
	return nil
}
