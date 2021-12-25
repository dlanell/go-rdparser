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
	text            string
	expectedProgram *Program
	expectedError   error
}

func TestRun(t *testing.T) {
	t.Run("Program", func(t *testing.T) {
		t.Run("NumericLiteral", func(t *testing.T) {
			tests := map[string]test{
				"given numbers": {
					text: `123;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{{
							nodeType: ExpressionStatement,
							body: &Node{
								nodeType: NumericLiteral,
								body: &NumericLiteralValue{
									value: 123,
								}},
						}},
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					parser := New(Props{Text: tc.text})
					node, err := parser.Run()
					assert.Equal(t, tc.expectedProgram, node)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
		t.Run("StringLiteral", func(t *testing.T) {
			tests := map[string]test{
				"given double quote string": {
					text: `"sith";`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{{
							nodeType: ExpressionStatement,
							body: &Node{
								nodeType: StringLiteral,
								body: &StringLiteralValue{
									value: "sith",
								}},
						}},
					},
				},
				"given single quote string": {
					text: `'sith';`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{{
							nodeType: ExpressionStatement,
							body: &Node{
								nodeType: StringLiteral,
								body: &StringLiteralValue{
									value: "sith",
								},
							},
						}},
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					parser := New(Props{Text: tc.text})
					node, err := parser.Run()
					assert.Equal(t, tc.expectedProgram, node)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
		t.Run("ExpressionStatement", func(t *testing.T) {
			tests := map[string]test{
				"given string and numeric expression": {
					text: `
'sith';
42;
`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: StringLiteral,
									body: &StringLiteralValue{
										value: "sith",
									},
								},
							},
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: NumericLiteral,
									body: &NumericLiteralValue{
										value: 42,
									},
								},
							},
						},
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					parser := New(Props{Text: tc.text})
					node, err := parser.Run()
					assert.Equal(t, tc.expectedProgram, node)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
		t.Run("BlockStatement", func(t *testing.T) {
			tests := map[string]test{
				"given block statement with expressions": {
					text: `
{
  'sith';
  42;
}
`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: BlockStatement,
								body: []*Node{
									{
										nodeType: ExpressionStatement,
										body: &Node{
											nodeType: StringLiteral,
											body:     &StringLiteralValue{"sith"},
										},
									},
									{
										nodeType: ExpressionStatement,
										body: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{42},
										},
									},
								},
							},
						},
					},
				},
				"given nested block statements with expressions": {
					text: `
{
  'sith';
  {
	42;
  }
}
`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: BlockStatement,
								body: []*Node{
									{
										nodeType: ExpressionStatement,
										body: &Node{
											nodeType: StringLiteral,
											body:     &StringLiteralValue{"sith"},
										},
									},
									{
										nodeType: BlockStatement,
										body: []*Node{{
											nodeType: ExpressionStatement,
											body:     &Node{
												nodeType: NumericLiteral,
												body: &NumericLiteralValue{42},
											},
										}},
									},
								},
							},
						},
					},
				},
				"given block statement without expressions": {
					text: `{}`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: BlockStatement,
								body: []*Node{},
							},
						},
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					parser := New(Props{Text: tc.text})
					node, err := parser.Run()
					assert.Equal(t, tc.expectedProgram, node)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
		t.Run("EmptyStatement", func(t *testing.T) {
			tests := map[string]test{
				"given empty statement": {
					text: `;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: EmptyStatement,
								body: nil,
							},
						},
					},
				},
			}

			for name, tc := range tests {
				t.Run(name, func(t *testing.T) {
					parser := New(Props{Text: tc.text})
					node, err := parser.Run()
					assert.Equal(t, tc.expectedProgram, node)
					assert.Equal(t, tc.expectedError, err)
				})
			}
		})
	})
}
