package tokenizer

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
			"given ;": {
				tokenizerText: `;`,
				expectedToken: &Token{
					TokenType: SemiColonToken,
					Value:     ";",
				},
			},
			"given {": {
				tokenizerText: `{`,
				expectedToken: &Token{
					TokenType: OpenCurlyBrace,
					Value:     "{",
				},
			},
			"given }": {
				tokenizerText: `}`,
				expectedToken: &Token{
					TokenType: CloseCurlyBrace,
					Value:     "}",
				},
			},
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
	t.Run("Math operators", func(t *testing.T) {
		tests := map[string]test{
			"given +": {
				tokenizerText: `+`,
				expectedToken: &Token{
					TokenType: AdditiveOperator,
					Value:     "+",
				},
			},
			"given -": {
				tokenizerText: `-`,
				expectedToken: &Token{
					TokenType: AdditiveOperator,
					Value:     "-",
				},
			},
			"given *": {
				tokenizerText: `*`,
				expectedToken: &Token{
					TokenType: MultiplicativeOperator,
					Value:     "*",
				},
			},
			"given /": {
				tokenizerText: `/`,
				expectedToken: &Token{
					TokenType: MultiplicativeOperator,
					Value:     "/",
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
	t.Run("Comment", func(t *testing.T) {
		tests := map[string]test{
			"given number after single line comment": {
				tokenizerText: `
// comment
1
`,
				expectedToken: &Token{
					TokenType: NumberToken,
					Value:     "1",
				},
			},
			"given number after multi line comment": {
				tokenizerText: `
/* 
comment
*/
1
`,
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
	t.Run("Assignment", func(t *testing.T) {
		tests := map[string]test{
			"given =": {
				tokenizerText: `=`,
				expectedToken: &Token{
					TokenType: SimpleAssignment,
					Value:     `=`,
				},
			},
			"given +=": {
				tokenizerText: `+=`,
				expectedToken: &Token{
					TokenType: ComplexAssignment,
					Value:     `+=`,
				},
			},
			"given -=": {
				tokenizerText: `-=`,
				expectedToken: &Token{
					TokenType: ComplexAssignment,
					Value:     `-=`,
				},
			},
			"given *=": {
				tokenizerText: `*=`,
				expectedToken: &Token{
					TokenType: ComplexAssignment,
					Value:     `*=`,
				},
			},
			"given /=": {
				tokenizerText: `/=`,
				expectedToken: &Token{
					TokenType: ComplexAssignment,
					Value:     `/=`,
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
	t.Run("Keywords", func(t *testing.T) {
		tests := map[string]test{
			"given let": {
				tokenizerText: `let`,
				expectedToken: &Token{
					TokenType: LetKeyword,
					Value:     `let`,
				},
			},
			"given if": {
				tokenizerText: `if`,
				expectedToken: &Token{
					TokenType: IfKeyword,
					Value:     `if`,
				},
			},
			"given else": {
				tokenizerText: `else`,
				expectedToken: &Token{
					TokenType: ElseKeyword,
					Value:     `else`,
				},
			},
			"given true": {
				tokenizerText: `true`,
				expectedToken: &Token{
					TokenType: TrueKeyword,
					Value:     `true`,
				},
			},
			"given false": {
				tokenizerText: `false`,
				expectedToken: &Token{
					TokenType: FalseKeyword,
					Value:     `false`,
				},
			},
			"given null": {
				tokenizerText: `null`,
				expectedToken: &Token{
					TokenType: NullKeyword,
					Value:     `null`,
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
	t.Run("Relational Operator", func(t *testing.T) {
		tests := map[string]test{
			"given >": {
				tokenizerText: `>`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `>`,
				},
			},
			"given >=": {
				tokenizerText: `>=`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `>=`,
				},
			},
			"given <": {
				tokenizerText: `<`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `<`,
				},
			},
			"given <=": {
				tokenizerText: `<=`,
				expectedToken: &Token{
					TokenType: RelationalOperator,
					Value:     `<=`,
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
	t.Run("Equality Operator", func(t *testing.T) {
		tests := map[string]test{
			"given ==": {
				tokenizerText: `==`,
				expectedToken: &Token{
					TokenType: EqualityOperator,
					Value:     `==`,
				},
			},
			"given !=": {
				tokenizerText: `!=`,
				expectedToken: &Token{
					TokenType: EqualityOperator,
					Value:     `!=`,
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
