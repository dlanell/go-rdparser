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
	body     interface{}
}

type Node struct {
	nodeType string
	value    interface{}
}

const (
	NumericLiteral string = "NUMERIC_LITERAL"
	StringLiteral  = "STRING_LITERAL"
	ProgramEnum     = "PROGRAM"
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
//	: Literal
//	;
///*
func (p *Parser) Program() (*Program, error) {
	numericLiteral, err := p.Literal()
	if err != nil {
		return nil, err
	}
	return &Program{
		nodeType: ProgramEnum,
		body: numericLiteral,
	}, nil
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
	return nil, errors.New("literal: unexpected literal production")
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

	return &Node{nodeType: NumericLiteral, value: num}, nil
}

// StringLiteral
//	: STRING
///*
func (p *Parser) StringLiteral() (*Node, error) {
	token, tokenErr := p.eat(tokenizer.StringToken)
	if tokenErr != nil {
		return nil, tokenErr
	}


	return &Node{nodeType: StringLiteral, value: token.Value[1:len(token.Value)-1]}, nil
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
