package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assertion := assert.New(t)
	t.Run("given New with Props, return new Tokenizer", func(t *testing.T) {
		tokenizer := New(Props{Text: "hello"})
		assertion.Equal(&Tokenizer{
			text:   "hello",
			cursor: 0,
		}, tokenizer)
	})
}

type test struct {
	tokenizerText string
	expectation *Token
}

func TestGetNextToken(t *testing.T) {
	t.Run("NumberToken", func(t *testing.T) {
		tests := map[string]test{
			"given empty string": {
				tokenizerText: "",
				expectation: nil,
			},
			"given valid number": {
				tokenizerText: "123",
				expectation: &Token{
					TokenType: NumberToken,
					Value:     "123",
				},
			},
			"given non numeric characters": {
				tokenizerText: "abc",
				expectation: nil,
			},
			"given non numeric characters after number": {
				tokenizerText: "1abc",
				expectation: &Token{
					TokenType: NumberToken,
					Value:     "1",
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectation, token)
			})
		}
	})
	t.Run("StringToken", func(t *testing.T) {
		t.Run("double quote strings", func(t *testing.T) {
			tests := map[string]test{
				"given no value": {
					tokenizerText: ``,
					expectation: nil,
				},
				"given valid string": {
					tokenizerText: `"sith"`,
					expectation: &Token{
						TokenType: StringToken,
						Value:     `"sith"`,
					},
				},
				"given characters after end of string": {
					tokenizerText: `"sith"1`,
					expectation: &Token{
						TokenType: StringToken,
						Value:     `"sith"`,
					},
				},
				"given number string": {
					tokenizerText: `"123"`,
					expectation: &Token{
						TokenType: StringToken,
						Value:     `"123"`,
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tokenizer := New(Props{Text: tc.tokenizerText})
					token := tokenizer.GetNextToken()
					assert.Equal(t, tc.expectation, token)
				})
			}
		})
		t.Run("single quote strings", func(t *testing.T) {
			tests := map[string]test{
				"given no value": {
					tokenizerText: ``,
					expectation: nil,
				},
				"given valid string": {
					tokenizerText: `'sith'`,
					expectation: &Token{
						TokenType: StringToken,
						Value:     `'sith'`,
					},
				},
				"given characters after end of string": {
					tokenizerText: `'sith'1`,
					expectation: &Token{
						TokenType: StringToken,
						Value:     `'sith'`,
					},
				},
				"given number string": {
					tokenizerText: `'123'`,
					expectation: &Token{
						TokenType: StringToken,
						Value:     `'123'`,
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tokenizer := New(Props{Text: tc.tokenizerText})
					token := tokenizer.GetNextToken()
					assert.Equal(t, tc.expectation, token)
				})
			}
		})
	})
}
