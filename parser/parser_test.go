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
						NodeType: ProgramEnum,
						Body: []*Node{{
							NodeType: ExpressionStatement,
							Body: &Node{
								NodeType: NumericLiteral,
								Body: &NumericLiteralValue{
									Value: 123,
								}},
						}},
					},
				},
				"given double quote string": {
					text: `"sith";`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{{
							NodeType: ExpressionStatement,
							Body: &Node{
								NodeType: StringLiteral,
								Body: &StringLiteralValue{
									Value: "sith",
								}},
						}},
					},
				},
				"given single quote string": {
					text: `'sith';`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{{
							NodeType: ExpressionStatement,
							Body: &Node{
								NodeType: StringLiteral,
								Body: &StringLiteralValue{
									Value: "sith",
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: StringLiteral,
									Body: &StringLiteralValue{
										Value: "sith",
									},
								},
							},
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: NumericLiteral,
									Body: &NumericLiteralValue{
										Value: 42,
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{"x"},
										},
										right: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{"x"},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "*",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "-",
										left: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												operator: "+",
												left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{3},
												},
												right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
											},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{3},
										},
										right: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												operator: "-",
												left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
												right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "/",
										left: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												operator: "*",
												left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{3},
												},
												right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
											},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{3},
										},
										right: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												operator: "*",
												left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
												right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "*",
										left: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												operator: "+",
												left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{3},
												},
												right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
											},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
									},
								},
							},
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "+",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{35},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{24},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "*",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
									},
								},
							},
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body: &BinaryExpressionNode{
										operator: "/",
										left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{35},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{24},
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
					text:          `42 = 42;`,
					expectedError: errors.New("invalid left-hand side in assignment expression"),
				},
				"given x = 42": {
					text: `x = 42;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: AssignmentExpression,
									Body: &BinaryExpressionNode{
										operator: "=",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: AssignmentExpression,
									Body: &BinaryExpressionNode{
										operator: "+=",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: AssignmentExpression,
									Body: &BinaryExpressionNode{
										operator: "-=",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: AssignmentExpression,
									Body: &BinaryExpressionNode{
										operator: "*=",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: AssignmentExpression,
									Body: &BinaryExpressionNode{
										operator: "/=",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: AssignmentExpression,
									Body: &BinaryExpressionNode{
										operator: "=",
										left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										right: &Node{
											NodeType: AssignmentExpression,
											Body: &BinaryExpressionNode{
												operator: "=",
												left: &Node{
													NodeType: Identifier,
													Body:     &StringLiteralValue{`y`},
												},
												right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: BlockStatement,
								Body: []*Node{
									{
										NodeType: ExpressionStatement,
										Body: &Node{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"sith"},
										},
									},
									{
										NodeType: ExpressionStatement,
										Body: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: BlockStatement,
								Body: []*Node{
									{
										NodeType: ExpressionStatement,
										Body: &Node{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"sith"},
										},
									},
									{
										NodeType: BlockStatement,
										Body: []*Node{{
											NodeType: ExpressionStatement,
											Body: &Node{
												NodeType: NumericLiteral,
												Body:     &NumericLiteralValue{42},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: BlockStatement,
								Body:     []*Node{},
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
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: EmptyStatement,
								Body:     nil,
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
		t.Run("VariableStatement", func(t *testing.T) {
			tests := map[string]test{
				"given let x = 42": {
					text: `let x = 42;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: VariableStatement,
								Body: []*Node{
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											id:   &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											init: &Node{
												NodeType: NumericLiteral,
												Body:     &NumericLiteralValue{42},
											},
										},
									},
								},
							},
						},
					},
				},
				"given let x": {
					text: `let x;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: VariableStatement,
								Body: []*Node{
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											init: nil,
										},
									},
								},
							},
						},
					},
				},
				"given let x, y": {
					text: `let x, y;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: VariableStatement,
								Body: []*Node{
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											init: nil,
										},
									},
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`y`},
											},
											init: nil,
										},
									},
								},
							},
						},
					},
				},
				"given let x, y = 45": {
					text: `let x, y = 45;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: VariableStatement,
								Body: []*Node{
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											init: nil,
										},
									},
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`y`},
											},
											init: &Node{
												NodeType: NumericLiteral,
												Body:     &NumericLiteralValue{45},
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
	})
}
