package queryparser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dlanell/go-rdparser/queryparser/querytokenizer"
)

/*

Program
	: Expression

Expression
	: Function | Literal

Function
	: LogicalFunc | RelationalFunction

LogicalFunction
	: logicalOperator '(' [Arguments] ')'

logicalOperator
	: and | or

Arguments
	: Expression [',' Arguments]

RelationalFunction
	: RelationalOperator '(' Identifier, Literal ')'

RelationalOperator
	: eq | ne | gt | ge | lt | le

Identifier
	: IDENTIFIER

Literal
	: NumericLiteral | StringLiteral | BooleanLiteral

NumericLiteral
	: NUMBER

StringLiteral
	: STRING

BooleanLiteral
	: true | false

*/

type QueryParser struct {
	text      string
	lookAhead *querytokenizer.Token
	tokenizer *querytokenizer.Tokenizer
}

type Props struct {
	Text string
}

type Program struct {
	NodeType string
	Body     *Node
}

type Node struct {
	NodeType string
	Body     interface{}
}

type FunctionNode struct {
	Operator  string
	Arguments []*Node
}

type BooleanLiteralValue struct {
	Value bool
}

type StringLiteralValue struct {
	Value string
}

type NumericLiteralValue struct {
	Value int
}

const (
	ProgramEnum        string = "Program"
	Identifier                = "Identifier"
	NumericLiteral            = "NumericLiteral"
	StringLiteral             = "StringLiteral"
	BooleanLiteral             = "BooleanLiteral"
	LogicalFunction           = "LogicalFunction"
	RelationalFunction        = "RelationalFunction"
)

func New(props Props) *QueryParser {
	return &QueryParser{
		text:      props.Text,
		tokenizer: querytokenizer.New(querytokenizer.Props{Text: props.Text}),
		lookAhead: nil,
	}
}

func (q *QueryParser) Run() (*Program, error) {
	token, err := q.tokenizer.GetNextToken()
	if err != nil {
		return nil, err
	}
	q.lookAhead = token

	return q.Program()
}

// Program
// Main entry point
//
// Program
//	: Expression
//	;
///*
func (q *QueryParser) Program() (*Program, error) {
	expression, err := q.Expression()
	if err != nil {
		return nil, err
	}
	return &Program{
		NodeType: ProgramEnum,
		Body:     expression,
	}, nil
}

// Expression
//	: Function | Literal
//	;
///*
func (q *QueryParser) Expression() (*Node, error) {
	switch q.lookAhead.TokenType {
	case querytokenizer.LogicalOperator:
		return q.Function()
	case querytokenizer.RelationalOperator:
		return q.Function()
	default:
		return q.Literal()
	}
}

// Function
//	: LogicalFunction | RelationalFunction
//	;
///*
func (q *QueryParser) Function() (*Node, error) {
	switch q.lookAhead.TokenType {
	case querytokenizer.LogicalOperator:
		return q.LogicalFunction()
	case querytokenizer.RelationalOperator:
		return q.RelationalFunction()
	}
	return nil, fmt.Errorf("Unexpected token: %s\n", q.lookAhead.TokenType)
}

// RelationalFunction
//	: RelationalOperator '(' Identifier, Literal ')'
//	;
///*
func (q *QueryParser) RelationalFunction() (*Node, error) {
	operator, err := q.eat(querytokenizer.RelationalOperator)
	if err != nil {
		return nil, err
	}
	_, err = q.eat(querytokenizer.OpenParentheses)
	if err != nil {
		return nil, err
	}

	identifier, firstArgumentErr := q.Identifier()
	if firstArgumentErr != nil {
		return nil, firstArgumentErr
	}

	_, err = q.eat(querytokenizer.Comma)
	if err != nil {
		return nil, err
	}

	literal, secondArgumentErr := q.Literal()
	if secondArgumentErr != nil {
		return nil, secondArgumentErr
	}

	_, err = q.eat(querytokenizer.CloseParentheses)
	if err != nil {
		return nil, err
	}
	arguments := make([]*Node, 0)
	arguments = append(arguments, identifier)
	arguments = append(arguments, literal)

	return &Node{
		NodeType: RelationalFunction,
		Body: &FunctionNode{
			Operator:  operator.Value,
			Arguments: arguments,
		},
	}, nil
}

// LogicalFunction
//	: logicalOperator '(' [Arguments] ')'
//	;
///*
func (q *QueryParser) LogicalFunction() (*Node, error) {
	operator, err := q.eat(querytokenizer.LogicalOperator)
	if err != nil {
		return nil, err
	}
	_, err = q.eat(querytokenizer.OpenParentheses)
	if err != nil {
		return nil, err
	}
	arguments, argumentsErr := q.Arguments()
	if argumentsErr != nil {
		return nil, argumentsErr
	}

	_, err = q.eat(querytokenizer.CloseParentheses)
	if err != nil {
		return nil, err
	}

	return &Node{
		NodeType: LogicalFunction,
		Body: &FunctionNode{
			Operator:  operator.Value,
			Arguments: arguments,
		},
	}, nil
}

// Arguments
//	: Expression [',' Arguments]
//	;
///*
func (q *QueryParser) Arguments() ([]*Node, error) {
	arguments := make([]*Node, 0)

	firstExpression, err := q.Expression()
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, firstExpression)

	_, err = q.eat(querytokenizer.Comma)
	if err != nil {
		return nil, err
	}

	followingExpression, followingExpressionErr := q.Expression()
	if followingExpressionErr != nil {
		return nil, followingExpressionErr
	}
	arguments = append(arguments, followingExpression)

	for q.lookAhead != nil && q.lookAhead.TokenType == querytokenizer.Comma {
		_, err = q.eat(querytokenizer.Comma)
		if err != nil {
			return nil, err
		}

		followingExpression, followingExpressionErr = q.Expression()
		if followingExpressionErr != nil {
			return nil, followingExpressionErr
		}
		arguments = append(arguments, followingExpression)
	}

	return arguments, nil
}

// Literal
//	: NumericLiteral | StringLiteral | BooleanLiteral
//	;
///*
func (q *QueryParser) Literal() (*Node, error) {
	switch q.lookAhead.TokenType {
	case querytokenizer.BooleanToken:
		return q.BooleanLiteral()
	case querytokenizer.NumberToken:
		return q.NumericLiteral()
	case querytokenizer.StringToken:
		return q.StringLiteral()
	}
	return nil, fmt.Errorf("unexpected token: %s\n", q.lookAhead.TokenType)
}

// NumericLiteral
//	: NUMBER
///*
func (q *QueryParser) NumericLiteral() (*Node, error) {
	token, tokenErr := q.eat(querytokenizer.NumberToken)
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
func (q *QueryParser) StringLiteral() (*Node, error) {
	token, tokenErr := q.eat(querytokenizer.StringToken)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return &Node{NodeType: StringLiteral, Body: &StringLiteralValue{token.Value[1 : len(token.Value)-1]}}, nil
}

// BooleanLiteral
//	: true | false
///*
func (q *QueryParser) BooleanLiteral() (*Node, error) {
	token, tokenErr := q.eat(querytokenizer.BooleanToken)
	if tokenErr != nil {
		return nil, tokenErr
	}

	value, valueErr := strconv.ParseBool(token.Value)
	if valueErr != nil {
		return nil, valueErr
	}

	return &Node{NodeType: BooleanLiteral, Body: &BooleanLiteralValue{value}}, nil
}

// Identifier
//	: IDENTIFIER
///*
func (q *QueryParser) Identifier() (*Node, error) {
	token, tokenErr := q.eat(querytokenizer.Identifier)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return &Node{NodeType: Identifier, Body: &StringLiteralValue{token.Value}}, nil
}

func (q *QueryParser) eat(tokenType string) (*querytokenizer.Token, error) {
	token := q.lookAhead
	if token == nil {
		return nil, fmt.Errorf("unexpected end of input, expected: %s\n", tokenType)
	}

	if token.TokenType != tokenType {
		return nil, fmt.Errorf("unexpected token: %s, expected: %s\n", token.Value, tokenType)
	}

	nextToken, _ := q.tokenizer.GetNextToken()
	q.lookAhead = nextToken

	return token, nil
}
