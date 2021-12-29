package queryparser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dlanell/go-rdparser/queryparser/querytokenizer"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assertion := assert.New(t)
	t.Run("given New with Props, return new QueryParser", func(t *testing.T) {
		parser := New(Props{Text: "hello"})
		assertion.Equal(&QueryParser{
			text:      "hello",
			tokenizer: querytokenizer.New(querytokenizer.Props{Text: "hello"}),
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
		t.Run("Number Literals", func(t *testing.T) {
			tests := map[string]test{
				"given numbers": {
					text: `123`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: NumericLiteral,
							Body: &NumericLiteralValue{
								Value: 123,
							},
						},
					},
				},
				"given invalid characters": {
					text:          `+`,
					expectedError: fmt.Errorf("unexpected token: %s", `+`),
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
		t.Run("Boolean Literals", func(t *testing.T) {
			tests := map[string]test{
				"given true": {
					text: `true`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: BooleanLiteral,
							Body: &BooleanLiteralValue{true},
						},
					},
				},
				"given false": {
					text: `false`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: BooleanLiteral,
							Body: &BooleanLiteralValue{false},
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
		t.Run("String Literals", func(t *testing.T) {
			tests := map[string]test{
				"given double quote string": {
					text: `"sith"`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: StringLiteral,
							Body: &StringLiteralValue{
								Value: "sith",
							},
						},
					},
				},
				"given single quote string": {
					text: `'sith'`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: StringLiteral,
							Body: &StringLiteralValue{
								Value: "sith",
							},
						},
					},
				},
				"given number string": {
					text: `'123'`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: StringLiteral,
							Body: &StringLiteralValue{
								Value: "123",
							},
						},
					},
				},
				"given invalid characters": {
					text:          `+`,
					expectedError: fmt.Errorf("unexpected token: %s", `+`),
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
		t.Run("Logical Function", func(t *testing.T) {
			t.Run("and", func(t *testing.T) {
				tests := map[string]test{
					"given multiple string literals": {
						text: `and("sith", "revan")`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "and",
									Arguments: []*Node{
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"sith"},
										},
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"revan"},
										},
									},
								},
							},
						},
					},
					"given multiple number literals": {
						text: `and(1, 2)`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "and",
									Arguments: []*Node{
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{1},
										},
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
					"given number & string literals": {
						text: `and(35, "revan")`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "and",
									Arguments: []*Node{
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{35},
										},
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"revan"},
										},
									},
								},
							},
						},
					},
					"given literals & functions": {
						text: `and(or(eq(path, "sith"), eq(name, "revan")), lt(count, 2), "star wars", 3, true)`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "and",
									Arguments: []*Node{
										{
											NodeType: LogicalFunction,
											Body: &FunctionNode{
												Operator: "or",
												Arguments: []*Node{
													{
														NodeType: RelationalFunction,
														Body: &FunctionNode{
															Operator: "eq",
															Arguments: []*Node{
																{
																	NodeType: Identifier,
																	Body:     &StringLiteralValue{"path"},
																},
																{
																	NodeType: StringLiteral,
																	Body:     &StringLiteralValue{"sith"},
																},
															},
														},
													},
													{
														NodeType: RelationalFunction,
														Body: &FunctionNode{
															Operator: "eq",
															Arguments: []*Node{
																{
																	NodeType: Identifier,
																	Body:     &StringLiteralValue{"name"},
																},
																{
																	NodeType: StringLiteral,
																	Body:     &StringLiteralValue{"revan"},
																},
															},
														},
													},
												},
											},
										},
										{
											NodeType: RelationalFunction,
											Body: &FunctionNode{
												Operator: "lt",
												Arguments: []*Node{
													{
														NodeType: Identifier,
														Body:     &StringLiteralValue{"count"},
													},
													{
														NodeType: NumericLiteral,
														Body:     &NumericLiteralValue{2},
													},
												},
											},
										},
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"star wars"},
										},
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{3},
										},
										{
											NodeType: BooleanLiteral,
											Body:     &BooleanLiteralValue{true},
										},
									},
								},
							},
						},
					},
					"given single expression": {
						text:          `and(35)`,
						expectedError: errors.New("unexpected token: ), expected: ,\n"),
					},
					"given single expression w/ no close parentheses": {
						text:          `and(1`,
						expectedError: errors.New("unexpected end of input, expected: ,\n"),
					},
					"given expressions w/ no close parentheses": {
						text:          `and(1, 3`,
						expectedError: errors.New("unexpected end of input, expected: )\n"),
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
			t.Run("or", func(t *testing.T) {
				tests := map[string]test{
					"given multiple string literals": {
						text: `or("sith", "revan")`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "or",
									Arguments: []*Node{
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"sith"},
										},
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"revan"},
										},
									},
								},
							},
						},
					},
					"given multiple number literals": {
						text: `or(1, 2)`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "or",
									Arguments: []*Node{
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{1},
										},
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{2},
										},
									},
								},
							},
						},
					},
					"given number & string literals": {
						text: `or(35, "revan")`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "or",
									Arguments: []*Node{
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{35},
										},
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"revan"},
										},
									},
								},
							},
						},
					},
					"given literals & functions": {
						text: `or(and(eq(path, "sith"), eq(name, "revan")), lt(count, 2), "star wars", 3, true)`,
						expectedProgram: &Program{
							NodeType: ProgramEnum,
							Body: &Node{
								NodeType: LogicalFunction,
								Body: &FunctionNode{
									Operator: "or",
									Arguments: []*Node{
										{
											NodeType: LogicalFunction,
											Body: &FunctionNode{
												Operator: "and",
												Arguments: []*Node{
													{
														NodeType: RelationalFunction,
														Body: &FunctionNode{
															Operator: "eq",
															Arguments: []*Node{
																{
																	NodeType: Identifier,
																	Body:     &StringLiteralValue{"path"},
																},
																{
																	NodeType: StringLiteral,
																	Body:     &StringLiteralValue{"sith"},
																},
															},
														},
													},
													{
														NodeType: RelationalFunction,
														Body: &FunctionNode{
															Operator: "eq",
															Arguments: []*Node{
																{
																	NodeType: Identifier,
																	Body:     &StringLiteralValue{"name"},
																},
																{
																	NodeType: StringLiteral,
																	Body:     &StringLiteralValue{"revan"},
																},
															},
														},
													},
												},
											},
										},
										{
											NodeType: RelationalFunction,
											Body: &FunctionNode{
												Operator: "lt",
												Arguments: []*Node{
													{
														NodeType: Identifier,
														Body:     &StringLiteralValue{"count"},
													},
													{
														NodeType: NumericLiteral,
														Body:     &NumericLiteralValue{2},
													},
												},
											},
										},
										{
											NodeType: StringLiteral,
											Body:     &StringLiteralValue{"star wars"},
										},
										{
											NodeType: NumericLiteral,
											Body:     &NumericLiteralValue{3},
										},
										{
											NodeType: BooleanLiteral,
											Body:     &BooleanLiteralValue{true},
										},
									},
								},
							},
						},
					},
					"given single expression": {
						text:          `or(35)`,
						expectedError: errors.New("unexpected token: ), expected: ,\n"),
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
		t.Run("Relational Function", func(t *testing.T) {
			tests := map[string]test{
				"given eq function": {
					text: `eq(sith, "revan")`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: RelationalFunction,
							Body: &FunctionNode{
								Operator: "eq",
								Arguments: []*Node{
									{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"sith"},
									},
									{
										NodeType: StringLiteral,
										Body:     &StringLiteralValue{"revan"},
									},
								},
							},
						},
					},
				},
				"given ne function": {
					text: `ne(sith, "revan")`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: RelationalFunction,
							Body: &FunctionNode{
								Operator: "ne",
								Arguments: []*Node{
									{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"sith"},
									},
									{
										NodeType: StringLiteral,
										Body:     &StringLiteralValue{"revan"},
									},
								},
							},
						},
					},
				},
				"given gt function": {
					text: `gt(count, 1)`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: RelationalFunction,
							Body: &FunctionNode{
								Operator: "gt",
								Arguments: []*Node{
									{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"count"},
									},
									{
										NodeType: NumericLiteral,
										Body:     &NumericLiteralValue{1},
									},
								},
							},
						},
					},
				},
				"given ge function": {
					text: `ge(count, 1)`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: RelationalFunction,
							Body: &FunctionNode{
								Operator: "ge",
								Arguments: []*Node{
									{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"count"},
									},
									{
										NodeType: NumericLiteral,
										Body:     &NumericLiteralValue{1},
									},
								},
							},
						},
					},
				},
				"given lt function": {
					text: `lt(count, 1)`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: RelationalFunction,
							Body: &FunctionNode{
								Operator: "lt",
								Arguments: []*Node{
									{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"count"},
									},
									{
										NodeType: NumericLiteral,
										Body:     &NumericLiteralValue{1},
									},
								},
							},
						},
					},
				},
				"given le function": {
					text: `le(count, 1)`,
					expectedProgram: &Program{
						NodeType: ProgramEnum,
						Body: &Node{
							NodeType: RelationalFunction,
							Body: &FunctionNode{
								Operator: "le",
								Arguments: []*Node{
									{
										NodeType: Identifier,
										Body:     &StringLiteralValue{"count"},
									},
									{
										NodeType: NumericLiteral,
										Body:     &NumericLiteralValue{1},
									},
								},
							},
						},
					},
				},
				"given only identifier": {
					text:          `le(count)`,
					expectedError: errors.New("unexpected token: ), expected: ,\n"),
				},
				"given invalid identifier": {
					text:          `le(42, 55)`,
					expectedError: errors.New("unexpected token: 42, expected: IDENTIFIER\n"),
				},
				"given no close parentheses": {
					text:          `le(count, 1`,
					expectedError: errors.New("unexpected end of input, expected: )\n"),
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
