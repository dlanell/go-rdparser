package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assertion := assert.New(t)
	t.Run("given New with Props, return new Tokenizer", func(t *testing.T) {
		tokenizer := New(Props{text: "hello"})
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
					tokenType: NumberToken,
					value:     "123",
				},
			},
			"given non numeric characters": {
				tokenizerText: "abc",
				expectation: nil,
			},
			"given non numeric characters after number": {
				tokenizerText: "1abc",
				expectation: &Token{
					tokenType: NumberToken,
					value:     "1",
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{text: tc.tokenizerText})
				token := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectation, token)
			})
		}
	})
}
