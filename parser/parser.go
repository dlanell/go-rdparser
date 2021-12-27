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
	NodeType string
	Body     []*Node
}

type Node struct {
	NodeType string
	Body     interface{}
}

type BinaryExpressionNode struct {
	Operator string
	Left     interface{}
	Right    interface{}
}

type VariableDeclarationValue struct {
	Id   *Node
	Init *Node
}

type IfStatementValue struct {
	Test   *Node
	Consequent *Node
	Alternate *Node
}

type StringLiteralValue struct {
	Value string
}

type NumericLiteralValue struct {
	Value int
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
	IfStatement                 = "IfStatement"
	VariableStatement           = "VariableStatement"
	VariableDeclaration         = "VariableDeclaration"
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
		NodeType: ProgramEnum,
		Body:     statements,
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
//	| EmptyStatement
//	| VariableStatement
//	| IfStatement
///*
func (p *Parser) Statement() (*Node, error) {
	switch p.lookAhead.TokenType {
	case tokenizer.SemiColonToken:
		return p.EmptyStatement()
	case tokenizer.OpenCurlyBrace:
		return p.BlockStatement()
	case tokenizer.LetKeyword:
		return p.VariableStatement()
	case tokenizer.IfKeyword:
		return p.IfStatement()
	default:
		return p.ExpressionStatement()
	}
}

// IfStatement
//	: 'if' '(' Expression ')' Statement
//	| 'if' '(' Expression ')' Statement 'else' Statement
///*
func (p *Parser) IfStatement() (*Node, error) {
	_, err := p.eat(tokenizer.IfKeyword)
	if err != nil {
		return nil, err
	}
	_, err = p.eat(tokenizer.OpenParentheses)
	if err != nil {
		return nil, err
	}

	test, testErr := p.Expression()
	if testErr != nil {
		return nil, testErr
	}

	_, err = p.eat(tokenizer.CloseParentheses)
	if err != nil {
		return nil, err
	}

	consequent, consequentErr := p.Statement()
	if consequentErr != nil {
		return nil, consequentErr
	}

	var alternate *Node
	var alternateErr error
	if p.lookAhead != nil && p.lookAhead.TokenType == tokenizer.ElseKeyword {
		_, err = p.eat(tokenizer.ElseKeyword)
		if err != nil {
			return nil, err
		}
		alternate, alternateErr = p.Statement()
		if alternateErr != nil {
			return nil, alternateErr
		}
	}

	return &Node{
		NodeType: IfStatement,
		Body:     &IfStatementValue{
			Test:       test,
			Consequent: consequent,
			Alternate:  alternate,
		},
	}, nil
}

// VariableStatement
//	: 'let' VariableDeclarationList ';'
///*
func (p *Parser) VariableStatement() (*Node, error) {
	_, err := p.eat(tokenizer.LetKeyword)
	if err != nil {
		return nil, err
	}
	declarationList, declarationListErr := p.VariableDeclarationList()
	if declarationListErr != nil {
		return nil, declarationListErr
	}

	_, err = p.eat(tokenizer.SemiColonToken)
	if err != nil {
		return nil, err
	}

	return &Node{NodeType: VariableStatement, Body: declarationList}, nil
}

// VariableDeclarationList
//	: VariableDeclarationList ',' VariableDeclaration
///*
func (p *Parser) VariableDeclarationList() ([]*Node, error) {
	declarations := make([]*Node, 0)

	for ok := true; ok; ok = p.lookAhead.TokenType == tokenizer.Comma {
		if p.lookAhead.TokenType == tokenizer.Comma {
			_, err := p.eat(tokenizer.Comma)
			if err != nil {
				return nil, err
			}
		}
		declaration, declarationErr := p.VariableDeclaration()
		if declarationErr != nil {
			return nil, declarationErr
		}
		declarations = append(declarations, declaration)
	}

	return declarations, nil
}

// VariableDeclaration
//	: Identifier OptVariableInitialization
///*
func (p *Parser) VariableDeclaration() (*Node, error) {
	identifier, err := p.Identifier()
	if err != nil {
		return nil, err
	}

	var init *Node
	var initErr error

	if p.lookAhead.TokenType != tokenizer.SemiColonToken && p.lookAhead.TokenType != tokenizer.Comma {
		init, initErr = p.VariableInitializer()
		if initErr != nil {
			return nil, initErr
		}
	}

	return &Node{
		NodeType: VariableDeclaration,
		Body: &VariableDeclarationValue{
			Id:   identifier,
			Init: init,
		},
	}, nil
}

// VariableInitializer
//	: SIMPLE_ASSIGNMENT AssignmentExpression
///*
func (p *Parser) VariableInitializer() (*Node, error) {
	_, err := p.eat(tokenizer.SimpleAssignment)
	if err != nil {
		return nil, err
	}
	return p.AssignmentExpression()
}

// EmptyStatement
//	: ';'
///*
func (p *Parser) EmptyStatement() (*Node, error) {
	_, err := p.eat(";")
	if err != nil {
		return nil, err
	}

	return &Node{NodeType: EmptyStatement, Body: nil}, nil
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
		return &Node{NodeType: BlockStatement, Body: []*Node{}}, nil
	}
	statements, statementsErr := p.StatementList(tokenizer.CloseCurlyBrace)
	if statementsErr != nil {
		return nil, err
	}
	_, err = p.eat(tokenizer.CloseCurlyBrace)
	if err != nil {
		return nil, err
	}

	return &Node{NodeType: BlockStatement, Body: statements}, nil
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

	return &Node{NodeType: ExpressionStatement, Body: expression}, nil
}

// Expression
//	: AssignmentExpression
///*
func (p *Parser) Expression() (*Node, error) {
	return p.AssignmentExpression()
}

// AssignmentExpression
//	: RelationalExpression
//	| LeftHandSideExpression AssignmentOperator AssignmentExpression
///*
func (p *Parser) AssignmentExpression() (*Node, error) {
	left, err := p.RelationalExpression()
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
		NodeType: AssignmentExpression,
		Body: &BinaryExpressionNode{
			Operator: assignmentOperatorToken.Value,
			Left:     leftNode,
			Right:    rightNode,
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
		NodeType: Identifier,
		Body:     &StringLiteralValue{token.Value},
	}, nil
}

func isAssignmentOperator(tokenType string) bool {
	return tokenType == tokenizer.SimpleAssignment || tokenType == tokenizer.ComplexAssignment
}

func checkValidAssignmentTarget(node *Node) (*Node, error) {
	if node.NodeType == Identifier {
		return node, nil
	}
	return nil, errors.New("invalid Left-hand side in assignment expression")
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

// RelationalExpression
//	: AdditiveExpression
//	| AdditiveExpression RELATIONAL_OPERATOR RelationalExpression
///*
func (p *Parser) RelationalExpression() (*Node, error) {
	return p.genericBinaryExpression(p.AdditiveExpression, tokenizer.RelationalOperator)
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
			NodeType: BinaryExpression,
			Body: &BinaryExpressionNode{
				Operator: operator.Value,
				Left:     left,
				Right:    right,
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

	return &Node{NodeType: NumericLiteral, Body: &NumericLiteralValue{Value: num}}, nil
}

// StringLiteral
//	: STRING
///*
func (p *Parser) StringLiteral() (*Node, error) {
	token, tokenErr := p.eat(tokenizer.StringToken)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return &Node{NodeType: StringLiteral, Body: &StringLiteralValue{token.Value[1 : len(token.Value)-1]}}, nil
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
