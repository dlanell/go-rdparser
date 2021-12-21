package parser

import (
	"testing"

	"github.com/dlanell/go-rdparser/parser/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assertion := assert.New(t)
	t.Run("given New with Props, return new Parser", func(t *testing.T) {
		parser := New(Props{Text: "hello"})
		assertion.Equal(&Parser{
			text:      "hello",
			tokenizer: tokenizer.New(tokenizer.Props{Text: "hello"}),
		}, parser)
	})
}

type test struct {
	tokenizerText string
	expectedNode  *Program
	expectedError error
}

func TestRun(t *testing.T) {
	t.Run("Program", func(t *testing.T) {
		tests := map[string]test{
			"given numbers": {
				tokenizerText: "123",
				expectedNode: &Program{
					nodeType: ProgramEnum,
					body: &Node{
						nodeType: NumericLiteral,
						value:    123,
					},
				},
			},
			"given double quote string": {
				tokenizerText: `"sith"`,
				expectedNode: &Program{
					nodeType: ProgramEnum,
					body: &Node{
						nodeType: StringLiteral,
						value:    "sith",
					},
				},
			},
			"given single quote string": {
				tokenizerText: `'sith'`,
				expectedNode: &Program{
					nodeType: ProgramEnum,
					body: &Node{
						nodeType: StringLiteral,
						value:    "sith",
					},
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				parser := New(Props{Text: tc.tokenizerText})
				node, err := parser.Run()
				assert.Equal(t, tc.expectedNode, node)
				assert.Equal(t, tc.expectedError, err)
			})
		}
	})
}
