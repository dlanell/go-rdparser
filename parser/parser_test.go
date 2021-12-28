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
		t.Run("Literals", func(t *testing.T) {
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
				"given true keyword": {
					text: `true;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{{
							NodeType: ExpressionStatement,
							Body: &Node{
								NodeType: BooleanLiteral,
								Body: &StringLiteralValue{
									Value: "true",
								},
							},
						}},
					},
				},
				"given false keyword": {
					text: `false;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{{
							NodeType: ExpressionStatement,
							Body: &Node{
								NodeType: BooleanLiteral,
								Body: &StringLiteralValue{
									Value: "false",
								},
							},
						}},
					},
				},
				"given null keyword": {
					text: `null;`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{{
							NodeType: ExpressionStatement,
							Body: &Node{
								NodeType: NullLiteral,
								Body: &StringLiteralValue{
									Value: "null",
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
										Operator: "+",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										Right: &Node{
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
										Operator: "+",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{"x"},
										},
										Right: &Node{
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
										Operator: "*",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										Right: &Node{
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
										Operator: "-",
										Left: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												Operator: "+",
												Left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{3},
												},
												Right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
											},
										},
										Right: &Node{
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
										Operator: "+",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{3},
										},
										Right: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												Operator: "-",
												Left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
												Right: &Node{
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
										Operator: "/",
										Left: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												Operator: "*",
												Left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{3},
												},
												Right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
											},
										},
										Right: &Node{
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
										Operator: "+",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{3},
										},
										Right: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												Operator: "*",
												Left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
												Right: &Node{
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
										Operator: "*",
										Left: &Node{
											NodeType: BinaryExpression,
											Body: &BinaryExpressionNode{
												Operator: "+",
												Left: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{3},
												},
												Right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{2},
												},
											},
										},
										Right: &Node{
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
										Operator: "+",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										Right: &Node{
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
										Operator: "+",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{35},
										},
										Right: &Node{
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
										Operator: "*",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
										Right: &Node{
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
										Operator: "/",
										Left: &Node{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{35},
										},
										Right: &Node{
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
					expectedError: errors.New("invalid Left-hand side in assignment expression"),
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
										Operator: "=",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										Right: &Node{
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
										Operator: "+=",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										Right: &Node{
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
										Operator: "-=",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										Right: &Node{
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
										Operator: "*=",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										Right: &Node{
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
										Operator: "/=",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										Right: &Node{
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
										Operator: "=",
										Left: &Node{
											NodeType: Identifier,
											Body:     &StringLiteralValue{`x`},
										},
										Right: &Node{
											NodeType: AssignmentExpression,
											Body: &BinaryExpressionNode{
												Operator: "=",
												Left: &Node{
													NodeType: Identifier,
													Body:     &StringLiteralValue{`y`},
												},
												Right: &Node{
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
											Id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											Init: &Node{
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
											Id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											Init: nil,
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
											Id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											Init: nil,
										},
									},
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											Id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`y`},
											},
											Init: nil,
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
											Id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`x`},
											},
											Init: nil,
										},
									},
									{
										NodeType: VariableDeclaration,
										Body: &VariableDeclarationValue{
											Id: &Node{
												NodeType: Identifier,
												Body:     &StringLiteralValue{`y`},
											},
											Init: &Node{
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
		t.Run("IfStatement", func(t *testing.T) {
			tests := map[string]test{
				"given valid if else statement with literal as test": {
					text: `

if (x) {
  x = 1;
} else {
  x = 2;
}

`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: IfStatement,
								Body: &IfStatementValue{
									Test: &Node{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"x"},
									},
									Consequent: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
															NodeType: NumericLiteral,
															Body:     &NumericLiteralValue{1},
														},
													},
												},
											},
										},
									},
									Alternate: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
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
					},
				},
				"given valid if statement with literal as test": {
					text: `

if (x) {
  x = 1;
}

`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: IfStatement,
								Body: &IfStatementValue{
									Test: &Node{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"x"},
									},
									Consequent: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
															NodeType: NumericLiteral,
															Body:     &NumericLiteralValue{1},
														},
													},
												},
											},
										},
									},
									Alternate: nil,
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
		t.Run("RelationalExpression", func(t *testing.T) {
			tests := map[string]test{
				"given valid if else statement with x + 5 > 10 as test": {
					text: `

if (x + 5 > 10) {
  x = 1;
} else {
  x = 2;
}

`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: IfStatement,
								Body: &IfStatementValue{
									Test: &Node{
										NodeType: BinaryExpression,
										Body:     &BinaryExpressionNode{
											Operator: `>`,
											Left: &Node{
												NodeType: BinaryExpression,
												Body:     &BinaryExpressionNode{
													Operator: `+`,
													Left: &Node{
														NodeType: Identifier,
														Body:     &StringLiteralValue{`x`},
													},
													Right: &Node{
														NodeType: NumericLiteral,
														Body:     &NumericLiteralValue{5},
													},
												},
											},
											Right: &Node{
												NodeType: NumericLiteral,
												Body:     &NumericLiteralValue{10},
											},
										},
									},
									Consequent: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
															NodeType: NumericLiteral,
															Body:     &NumericLiteralValue{1},
														},
													},
												},
											},
										},
									},
									Alternate: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
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
		t.Run("EqualityExpression", func(t *testing.T) {
			tests := map[string]test{
				"given valid if else statement with x + 5 == 10 as test": {
					text: `

if (x + 5 == 10) {
  x = 1;
} else {
  x = 2;
}

`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: IfStatement,
								Body: &IfStatementValue{
									Test: &Node{
										NodeType: BinaryExpression,
										Body:     &BinaryExpressionNode{
											Operator: `==`,
											Left: &Node{
												NodeType: BinaryExpression,
												Body:     &BinaryExpressionNode{
													Operator: `+`,
													Left: &Node{
														NodeType: Identifier,
														Body:     &StringLiteralValue{`x`},
													},
													Right: &Node{
														NodeType: NumericLiteral,
														Body:     &NumericLiteralValue{5},
													},
												},
											},
											Right: &Node{
												NodeType: NumericLiteral,
												Body:     &NumericLiteralValue{10},
											},
										},
									},
									Consequent: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
															NodeType: NumericLiteral,
															Body:     &NumericLiteralValue{1},
														},
													},
												},
											},
										},
									},
									Alternate: &Node{
										NodeType: BlockStatement,
										Body: []*Node{
											{
												NodeType: ExpressionStatement,
												Body: &Node{
													NodeType: AssignmentExpression,
													Body:     &BinaryExpressionNode{
														Operator: "=",
														Left: &Node{
															NodeType: Identifier,
															Body:     &StringLiteralValue{"x"},
														},
														Right: &Node{
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
					},
				},
				"given x > 5 == true;": {
					text: `

x > 5 == true;

`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: []*Node{
							{
								NodeType: ExpressionStatement,
								Body: &Node{
									NodeType: BinaryExpression,
									Body:     &BinaryExpressionNode{
										Operator: `==`,
										Left: &Node{
											NodeType: BinaryExpression,
											Body:     &BinaryExpressionNode{
												Operator: `>`,
												Left: &Node{
													NodeType: Identifier,
													Body:     &StringLiteralValue{`x`},
												},
												Right: &Node{
													NodeType: NumericLiteral,
													Body:     &NumericLiteralValue{5},
												},
											},
										},
										Right: &Node{
											NodeType: BooleanLiteral,
											Body:     &StringLiteralValue{"true"},
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
