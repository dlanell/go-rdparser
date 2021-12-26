package parser

import (
	"errors"
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
		t.Run("ExpressionStatement", func(t *testing.T) {
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
		t.Run("Math BinaryExpression", func(t *testing.T) {
			tests := map[string]test{
				"given 2 + 2": {
					text: `2 + 2;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
				},
				"given x + x": {
					text: `x + x;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{"x"},
										},
										right: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{"x"},
										},
									},
								},
							},
						},
					},
				},
				"given 2 * 2": {
					text: `2 * 2;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "*",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
				},
				"given chained additive operators 3 + 2 - 2": {
					text: `3 + 2 - 2;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "-",
										left: &Node{
											nodeType: BinaryExpression,
											body: &BinaryExpressionNode{
												operator: "+",
												left: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{3},
												},
												right: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
											},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
				},
				"given chained additive operators with parentheses 3 + ( 2 - 2 )": {
					text: `3 + ( 2 - 2 );`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{3},
										},
										right: &Node{
											nodeType: BinaryExpression,
											body: &BinaryExpressionNode{
												operator: "-",
												left: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
												right: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"given chained multiplicative operators 3 * 2 / 2": {
					text: `3 * 2 / 2;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "/",
										left: &Node{
											nodeType: BinaryExpression,
											body: &BinaryExpressionNode{
												operator: "*",
												left: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{3},
												},
												right: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
											},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
				},
				"given chained multiplicative & additive operators 3 + 2 * 2": {
					text: `3 + 2 * 2;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{3},
										},
										right: &Node{
											nodeType: BinaryExpression,
											body: &BinaryExpressionNode{
												operator: "*",
												left: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
												right: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"given chained multiplicative & additive operators with parentheses (3 + 2) * 2": {
					text: `(3 + 2) * 2;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "*",
										left: &Node{
											nodeType: BinaryExpression,
											body: &BinaryExpressionNode{
												operator: "+",
												left: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{3},
												},
												right: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{2},
												},
											},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
				},
				"given 2 additive expressions": {
					text: `
2 + 2;
35 + 24;
`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{35},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{24},
										},
									},
								},
							},
						},
					},
				},
				"given 2 multiplicative expressions": {
					text: `
2 * 2;
35 / 24;
`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "*",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{2},
										},
									},
								},
							},
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: BinaryExpression,
									body: &BinaryExpressionNode{
										operator: "/",
										left: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{35},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{24},
										},
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
		t.Run("Assignment BinaryExpression", func(t *testing.T) {
			tests := map[string]test{
				"given 42 = 42": {
					text: `42 = 42;`,
					expectedError: errors.New("invalid left-hand side in assignment expression"),
				},
				"given x = 42": {
					text: `x = 42;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: AssignmentExpression,
									body: &BinaryExpressionNode{
										operator: "=",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{42},
										},
									},
								},
							},
						},
					},
				},
				"given x += 42": {
					text: `x += 42;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: AssignmentExpression,
									body: &BinaryExpressionNode{
										operator: "+=",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{42},
										},
									},
								},
							},
						},
					},
				},
				"given x -= 42": {
					text: `x -= 42;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: AssignmentExpression,
									body: &BinaryExpressionNode{
										operator: "-=",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{42},
										},
									},
								},
							},
						},
					},
				},
				"given x *= 42": {
					text: `x *= 42;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: AssignmentExpression,
									body: &BinaryExpressionNode{
										operator: "*=",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{42},
										},
									},
								},
							},
						},
					},
				},
				"given x /= 42": {
					text: `x /= 42;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: AssignmentExpression,
									body: &BinaryExpressionNode{
										operator: "/=",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											nodeType: NumericLiteral,
											body:     &NumericLiteralValue{42},
										},
									},
								},
							},
						},
					},
				},
				"given x = y = 42": {
					text: `x = y = 42;`,
					expectedProgram: &Program{
						nodeType: ProgramEnum,
						body: []*Node{
							{
								nodeType: ExpressionStatement,
								body: &Node{
									nodeType: AssignmentExpression,
									body: &BinaryExpressionNode{
										operator: "=",
										left: &Node{
											nodeType: Identifier,
											body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											nodeType: AssignmentExpression,
											body: &BinaryExpressionNode{
												operator: "=",
												left: &Node{
													nodeType: Identifier,
													body:     &StringLiteralValue{`y`},
												},
												right: &Node{
													nodeType: NumericLiteral,
													body:     &NumericLiteralValue{42},
												},
											},
										},
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
											body: &Node{
												nodeType: NumericLiteral,
												body:     &NumericLiteralValue{42},
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
								body:     []*Node{},
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
								body:     nil,
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
