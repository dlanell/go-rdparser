package querytokenizer

import (
	"errors"
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
	expectedToken *Token
	expectedError error
}

func TestGetNextToken(t *testing.T) {
	t.Run("Symbols & Delimeters", func(t *testing.T) {
		tests := map[string]test{
			"given (": {
				tokenizerText: `(`,
				expectedToken: &Token{
					TokenType: OpenParentheses,
					Value:     "(",
				},
			},
			"given )": {
				tokenizerText: `)`,
				expectedToken: &Token{
					TokenType: CloseParentheses,
					Value:     ")",
				},
			},
			"given ,": {
				tokenizerText: `,`,
				expectedToken: &Token{
					TokenType: Comma,
					Value:     ",",
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
	t.Run("NumberToken", func(t *testing.T) {
		tests := map[string]test{
			"given empty string": {
				tokenizerText: ``,
				expectedError: errors.New("no tokens present"),
			},
			"given valid number": {
				tokenizerText: `123`,
				expectedToken: &Token{
					TokenType: NumberToken,
					Value:     `123`,
				},
			},
			"given valid number after whitespace": {
				tokenizerText: `        123`,
				expectedToken: &Token{
					TokenType: NumberToken,
					Value:     `123`,
				},
			},
			"given non numeric characters after number": {
				tokenizerText: `1a`,
				expectedToken: &Token{
					TokenType: NumberToken,
					Value:     "1",
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
	t.Run("DateToken", func(t *testing.T) {
		tests := map[string]test{
			"given date token": {
				tokenizerText: `2020-04-03T08:58:26Z`,
				expectedToken: &Token{
					TokenType: DateToken,
					Value:     `2020-04-03T08:58:26Z`,
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
	t.Run("BooleanToken", func(t *testing.T) {
		tests := map[string]test{
			"given empty string": {
				tokenizerText: ``,
				expectedError: errors.New("no tokens present"),
			},
			"given true": {
				tokenizerText: `true`,
				expectedToken: &Token{
					TokenType: BooleanToken,
					Value:     `true`,
				},
			},
			"given false": {
				tokenizerText: `    false    `,
				expectedToken: &Token{
					TokenType: BooleanToken,
					Value:     `false`,
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
	t.Run("StringToken", func(t *testing.T) {
		t.Run("double quote strings", func(t *testing.T) {
			tests := map[string]test{
				"given no value": {
					tokenizerText: ``,
					expectedError: errors.New("no tokens present"),
				},
				"given valid string": {
					tokenizerText: `"sith"`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `"sith"`,
					},
				},
				"given string with whitespace within quotes": {
					tokenizerText: `"  sith  "`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `"  sith  "`,
					},
				},
				"given valid string after whitespace": {
					tokenizerText: `        "sith"`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `"sith"`,
					},
				},
				"given characters after end of string": {
					tokenizerText: `"sith"1`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `"sith"`,
					},
				},
				"given number string": {
					tokenizerText: `"123"`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `"123"`,
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tokenizer := New(Props{Text: tc.tokenizerText})
					token, err := tokenizer.GetNextToken()
					assert.Equal(t, tc.expectedToken, token)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
		t.Run("single quote strings", func(t *testing.T) {
			tests := map[string]test{
				"given no value": {
					tokenizerText: ``,
					expectedError: errors.New("no tokens present"),
				},
				"given valid string": {
					tokenizerText: `'sith'`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `'sith'`,
					},
				},
				"given valid string with whitespace within quotes": {
					tokenizerText: `'  sith  '`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `'  sith  '`,
					},
				},
				"given valid string after whitespace": {
					tokenizerText: `      'sith'`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `'sith'`,
					},
				},
				"given characters after end of string": {
					tokenizerText: `'sith'1`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `'sith'`,
					},
				},
				"given number string": {
					tokenizerText: `'123'`,
					expectedToken: &Token{
						TokenType: StringToken,
						Value:     `'123'`,
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					tokenizer := New(Props{Text: tc.tokenizerText})
					token, err := tokenizer.GetNextToken()
					assert.Equal(t, tc.expectedToken, token)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
	})
	t.Run("Identifier", func(t *testing.T) {
		tests := map[string]test{
			"given windu": {
				tokenizerText: `windu`,
				expectedToken: &Token{
					TokenType: Identifier,
					Value:     `windu`,
				},
			},
			"given windu123": {
				tokenizerText: `windu123`,
				expectedToken: &Token{
					TokenType: Identifier,
					Value:     `windu123`,
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
	t.Run("Relational Operators", func(t *testing.T) {
		tests := map[string]test{
			"given eq": {
				tokenizerText: `eq`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `eq`,
				},
			},
			"given ne": {
				tokenizerText: `ne`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `ne`,
				},
			},
			"given lt": {
				tokenizerText: `lt`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `lt`,
				},
			},
			"given le": {
				tokenizerText: `le`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `le`,
				},
			},
			"given gt": {
				tokenizerText: `gt`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `gt`,
				},
			},
			"given ge": {
				tokenizerText: `ge`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `ge`,
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
	t.Run("Logical Operators", func(t *testing.T) {
		tests := map[string]test{
			"given and": {
				tokenizerText: `and`,
				expectedToken: &Token{
					TokenType: LogicalOperator,
					Value:     `and`,
				},
			},
			"given or": {
				tokenizerText: `or`,
				expectedToken: &Token{
					TokenType: LogicalOperator,
					Value:     `or`,
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				tokenizer := New(Props{Text: tc.tokenizerText})
				token, err := tokenizer.GetNextToken()
				assert.Equal(t, tc.expectedToken, token)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
}