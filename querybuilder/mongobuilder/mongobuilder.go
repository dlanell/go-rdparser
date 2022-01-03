package mongobuilder

import (
	"time"

	"github.com/dlanell/go-rdparser/queryparser"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoFilterBuilder struct {
	literalComparisonFields []string
	parseTree *queryparser.Node
	parser *queryparser.QueryParser
}

type Props struct {
	LiteralComparisonFields []string
}

var OperatorMap = map[string]string {
	queryparser.NotEqualOperator:           "$ne",
	queryparser.GreaterThanOperator:        "$gt",
	queryparser.GreaterThanOrEqualOperator: "$gte",
	queryparser.LessThanOperator:           "$lt",
	queryparser.LessThanOrEqualOperator:    "$lte",
	queryparser.LogicalOrOperator:          "$or",
	queryparser.LogicalAndOperator:         "$and",
}

func New(props Props) *MongoFilterBuilder {
	return &MongoFilterBuilder{
		literalComparisonFields: props.LiteralComparisonFields,
		parser: queryparser.New(),
		parseTree:      nil,
	}
}

func (m *MongoFilterBuilder) Run(text string) (bson.D,error) {
	tree, err := m.parser.Run(text)
	if err != nil {
		return nil, err
	}

	m.parseTree = tree.Body
	return m.parseNode(m.parseTree), nil
}

func (m *MongoFilterBuilder) parseNode(node *queryparser.Node) bson.D {
	switch node.NodeType {
	case queryparser.RelationalFunction:
		return m.parseRelationalFunctionNode(node.Body.(*queryparser.FunctionNode))
	case queryparser.StringLiteral:
		return m.parseLiteralNode(node)
	case queryparser.NumericLiteral:
		return m.parseLiteralNode(node)
	case queryparser.DateLiteral:
		return m.parseLiteralNode(node)
	case queryparser.BooleanLiteral:
		return m.parseLiteralNode(node)
	default:
		return m.parseLogicalFunctionNode(node.Body.(*queryparser.FunctionNode))
	}
}

func (m *MongoFilterBuilder) parseLiteralNode(node *queryparser.Node) bson.D {
	fieldComparisons := bson.A{}
	for _, field := range m.literalComparisonFields {
		fieldComparisons = append(fieldComparisons, bson.D{{field, m.getLiteralNodeValue(node)}})
	}
	return bson.D{{"$or", fieldComparisons}}
}

func (m *MongoFilterBuilder) parseLogicalFunctionNode(node *queryparser.FunctionNode) bson.D {
	argumentNodes := node.Arguments
	arguments := bson.A{}
	for _, argumentNode := range argumentNodes {
		arguments = append(arguments, m.parseNode(argumentNode))
	}

	return bson.D{{OperatorMap[node.Operator], arguments}}
}

func (m *MongoFilterBuilder) parseRelationalFunctionNode(node *queryparser.FunctionNode) bson.D {
	if node.Operator == queryparser.EqualOperator {
		return m.equalRelationalFunction(node)
	}
	return m.nonEqualRelationalFunction(node)
}

func (m *MongoFilterBuilder) equalRelationalFunction(node *queryparser.FunctionNode) bson.D  {
	arguments := node.Arguments
	identifier := m.getLiteralNodeValue(arguments[0]).(string)
	literal := m.getLiteralNodeValue(arguments[1])

	return bson.D{{identifier, literal}}
}

func (m *MongoFilterBuilder) nonEqualRelationalFunction(node *queryparser.FunctionNode) bson.D {
	arguments := node.Arguments
	identifier := m.getLiteralNodeValue(arguments[0]).(string)
	literal := m.getLiteralNodeValue(arguments[1])

	return bson.D{{identifier, bson.D{{OperatorMap[node.Operator], literal}}}}
}

func (m *MongoFilterBuilder) getLiteralNodeValue(node *queryparser.Node) interface{} {
	switch node.NodeType {
	case queryparser.BooleanLiteral:
		return node.Body.(*queryparser.BooleanLiteralValue).Value
	case queryparser.NumericLiteral:
		return node.Body.(*queryparser.NumericLiteralValue).Value
	case queryparser.Identifier:
		return node.Body.(*queryparser.StringLiteralValue).Value
	case queryparser.DateLiteral:
		utcTime, _ := time.Parse(time.RFC3339, node.Body.(*queryparser.StringLiteralValue).Value)
		return utcTime
	default:
		return node.Body.(*queryparser.StringLiteralValue).Value
	}

}