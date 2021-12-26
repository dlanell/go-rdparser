package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dlanell/go-rdparser/parser/tokenizer"
)

type Parser struct {
	text      string
	lookAhead *tokenizer.Token
	tokenizer *tokenizer.Tokenizer
}

type Props struct {
	Text string
}

type Program struct {
	nodeType string
	body     []*Node
}

type Node struct {
	nodeType string
	body     interface{}
}

type BinaryExpressionNode struct {
	operator string
	left     interface{}
	right    interface{}
}

type StringLiteralValue struct {
	value string
}

type NumericLiteralValue struct {
	value int
}

const (
	NumericLiteral       string = "NumericLiteral"
	StringLiteral               = "StringLiteral"
	Identifier                  = "IDENTIFIER"
	ExpressionStatement         = "ExpressionStatement"
	AssignmentExpression        = "AssignmentExpression"
	BlockStatement              = "BlockStatement"
	BinaryExpression            = "BinaryExpression"
	EmptyStatement              = "EmptyStatement"
	ProgramEnum                 = "Program"
)

func New(props Props) *Parser {
	return &Parser{
		text:      props.Text,
		tokenizer: tokenizer.New(tokenizer.Props{Text: props.Text}),
		lookAhead: nil,
	}
}

func (p *Parser) Run() (*Program, error) {
	token, err := p.tokenizer.GetNextToken()
	if err != nil {
		return nil, err
	}
	p.lookAhead = token

	return p.Program()
}

// Program
// Main entry point
//
// Program
//	: StatementList
//	;
///*
func (p *Parser) Program() (*Program, error) {
	statements, err := p.StatementList("")
	if err != nil {
		return nil, err
	}
	return &Program{
		nodeType: ProgramEnum,
		body:     statements,
	}, nil
}

// StatementList
//	: Statement
//	| Statement StatementList
///*
func (p *Parser) StatementList(stopLookAhead string) ([]*Node, error) {
	statements := make([]*Node, 0)
	var statement, err = p.Statement()
	if err != nil {
		return nil, err
	}
	statements = append(statements, statement)
	for p.lookAhead != nil && p.lookAhead.TokenType != stopLookAhead {
		statement, err = p.Statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

// Statement
//	: ExpressionStatement
//	| BlockStatement
///*
func (p *Parser) Statement() (*Node, error) {
	switch p.lookAhead.TokenType {
	case tokenizer.SemiColonToken:
		return p.EmptyStatement()
	case tokenizer.OpenCurlyBrace:
		return p.BlockStatement()
	default:
		return p.ExpressionStatement()
	}
}

// EmptyStatement
//	: ';'
///*
func (p *Parser) EmptyStatement() (*Node, error) {
	_, err := p.eat(";")
	if err != nil {
		return nil, err
	}

	return &Node{nodeType: EmptyStatement, body: nil}, nil
}

// BlockStatement
//	: '{' OptStatementList '}'
///*
func (p *Parser) BlockStatement() (*Node, error) {
	_, err := p.eat("{")
	if err != nil {
		return nil, err
	}
	if p.lookAhead.TokenType == tokenizer.CloseCurlyBrace {
		_, err = p.eat(tokenizer.CloseCurlyBrace)
		if err != nil {
			return nil, err
		}
		return &Node{nodeType: BlockStatement, body: []*Node{}}, nil
	}
	statements, statementsErr := p.StatementList(tokenizer.CloseCurlyBrace)
	if statementsErr != nil {
		return nil, err
	}
	_, err = p.eat(tokenizer.CloseCurlyBrace)
	if err != nil {
		return nil, err
	}

	return &Node{nodeType: BlockStatement, body: statements}, nil
}

// ExpressionStatement
//	: Expression ';'
///*
func (p *Parser) ExpressionStatement() (*Node, error) {
	expression, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.eat(";")
	if err != nil {
		return nil, err
	}

	return &Node{nodeType: ExpressionStatement, body: expression}, nil
}

// Expression
//	: AssignmentExpression
///*
func (p *Parser) Expression() (*Node, error) {
	return p.AssignmentExpression()
}

// AssignmentExpression
//	: AdditiveExpression
//	| LeftHandSideExpression AssignmentOperator AssignmentExpression
///*
func (p *Parser) AssignmentExpression() (*Node, error) {
	left, err := p.AdditiveExpression()
	if err != nil {
		return nil, err
	}

	if !isAssignmentOperator(p.lookAhead.TokenType) {
		return left, nil
	}

	assignmentOperatorToken, assignmentOperatorTokenErr := p.AssignmentOperator()
	if assignmentOperatorTokenErr != nil {
		return nil, assignmentOperatorTokenErr
	}

	leftNode, leftNodeErr := checkValidAssignmentTarget(left)
	if leftNodeErr != nil {
		return nil, leftNodeErr
	}

	rightNode, rightNodeErr := p.AssignmentExpression()
	if rightNodeErr != nil {
		return nil, rightNodeErr
	}

	return &Node{
		nodeType: AssignmentExpression,
		body: &BinaryExpressionNode{
			operator: assignmentOperatorToken.Value,
			left:     leftNode,
			right:    rightNode,
		},
	}, nil
}

// LeftHandSideExpression
//	: Identifier
///*
func (p *Parser) LeftHandSideExpression() (*Node, error) {
	return p.Identifier()
}

// Identifier
//	: IDENTIFIER
///*
func (p *Parser) Identifier() (*Node, error) {
	token, err := p.eat(tokenizer.Identifier)
	if err != nil {
		return nil, err
	}
	return &Node{
		nodeType: Identifier,
		body:     &StringLiteralValue{token.Value},
	}, nil
}

func isAssignmentOperator(tokenType string) bool {
	return tokenType == tokenizer.SimpleAssignment || tokenType == tokenizer.ComplexAssignment
}

func checkValidAssignmentTarget(node *Node) (*Node, error) {
	if node.nodeType == Identifier {
		return node, nil
	}
	return nil, errors.New("invalid left-hand side in assignment expression")
}

// AssignmentOperator
//	: Simple Assignment Token
//	| Complex Assignment Token
///*
func (p *Parser) AssignmentOperator() (*tokenizer.Token, error) {
	if p.lookAhead.TokenType == tokenizer.SimpleAssignment {
		return p.eat(tokenizer.SimpleAssignment)
	}
	return p.eat(tokenizer.ComplexAssignment)
}

// AdditiveExpression
//	: MultiplicativeExpression
//	| MultiplicativeExpression Additive_Operator AdditiveExpression
///*
func (p *Parser) AdditiveExpression() (*Node, error) {
	return p.genericBinaryExpression(p.MultiplicativeExpression, tokenizer.AdditiveOperator)
}

// MultiplicativeExpression
//	: PrimaryExpression
//	| PrimaryExpression MultiplicativeOperator MultiplicativeExpression
///*
func (p *Parser) MultiplicativeExpression() (*Node, error) {
	return p.genericBinaryExpression(p.PrimaryExpression, tokenizer.MultiplicativeOperator)
}

func (p *Parser) genericBinaryExpression(expression func() (*Node, error), operatorToken string) (*Node, error) {
	left, err := expression()
	if err != nil {
		return nil, err
	}

	for p.lookAhead.TokenType == operatorToken {
		operator, operatorErr := p.eat(operatorToken)
		if operatorErr != nil {
			return nil, operatorErr
		}
		right, rightErr := expression()
		if rightErr != nil {
			return nil, rightErr
		}

		left = &Node{
			nodeType: BinaryExpression,
			body: &BinaryExpressionNode{
				operator: operator.Value,
				left:     left,
				right:    right,
			},
		}
	}
	return left, nil
}

// PrimaryExpression
//	: Literal
//	| ParenthesizedExpression
//	| LeftHandSideExpression
///*
func (p *Parser) PrimaryExpression() (*Node, error) {
	if isLiteral(p.lookAhead.TokenType) {
		return p.Literal()
	}
	switch p.lookAhead.TokenType {
	case tokenizer.OpenParentheses:
		return p.ParenthesizedExpression()
	default:
		return p.LeftHandSideExpression()
	}
}

func isLiteral(tokenType string) bool {
	return tokenType == tokenizer.StringToken || tokenType == tokenizer.NumberToken
}

// ParenthesizedExpression
//	: '(' Expression ')'
//	;
///*
func (p *Parser) ParenthesizedExpression() (*Node, error) {
	_, err := p.eat(tokenizer.OpenParentheses)
	if err != nil {
		return nil, err
	}
	expression, expressionErr := p.Expression()
	if expressionErr != nil {
		return nil, expressionErr
	}
	_, err = p.eat(tokenizer.CloseParentheses)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

// Literal
//	: NumericLiteral
//	| StringLiteral
///*
func (p *Parser) Literal() (*Node, error) {
	switch p.lookAhead.TokenType {
	case tokenizer.NumberToken:
		return p.NumericLiteral()
	case tokenizer.StringToken:
		return p.StringLiteral()
	}
	return nil, fmt.Errorf("Unexpected token: %s\n", p.lookAhead.TokenType)
}

// NumericLiteral
//	: NUMBER
///*
func (p *Parser) NumericLiteral() (*Node, error) {
	token, tokenErr := p.eat(tokenizer.NumberToken)
	if tokenErr != nil {
		return nil, tokenErr
	}

	num, err := strconv.Atoi(token.Value)
	if err != nil {
		return nil, errors.New("invalid number token")
	}

	return &Node{nodeType: NumericLiteral, body: &NumericLiteralValue{value: num}}, nil
}

// StringLiteral
//	: STRING
///*
func (p *Parser) StringLiteral() (*Node, error) {
	token, tokenErr := p.eat(tokenizer.StringToken)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return &Node{nodeType: StringLiteral, body: &StringLiteralValue{token.Value[1 : len(token.Value)-1]}}, nil
}

func (p *Parser) eat(tokenType string) (*tokenizer.Token, error) {
	token := p.lookAhead
	if token == nil {
		return nil, fmt.Errorf("Unexpected end of input, expected: %s\n", tokenType)
	}

	if token.TokenType != tokenType {
		return nil, fmt.Errorf("Unexpected token: %s, expected: %s\n", token.Value, tokenType)
	}

	nextToken, _ := p.tokenizer.GetNextToken()
	p.lookAhead = nextToken

	return token, nil
}
