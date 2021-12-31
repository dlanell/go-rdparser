package mongobuilder

import (
	"testing"

	"github.com/dlanell/go-rdparser/queryparser"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type test struct {
	filterParam   string
	expectedQuery bson.D
	expectedError error
}

func TestRun(t *testing.T) {
	literalComparisonFields := []string{
		"title",
		"email",
	}
	t.Run("Relation Functions", func(t *testing.T) {
		tests := map[string]test{
			`given eq(policyId, "someId")`: {
				filterParam:   `eq(policyId, "someId")`,
				expectedQuery: bson.D{{"policyId", "someId"}},
			},
			`given eq(cores, 10)`: {
				filterParam:   `eq(cores, 10)`,
				expectedQuery: bson.D{{"cores", 10}},
			},
			`given ne(policyId, "someId")`: {
				filterParam:   `ne(policyId, "someId")`,
				expectedQuery: bson.D{{"policyId", bson.D{{"$ne", "someId"}}}},
			},
			`given ne(cores, 4)`: {
				filterParam:   `ne(cores, 4)`,
				expectedQuery: bson.D{{"cores", bson.D{{"$ne", 4}}}},
			},
			`given gt(policyId, "someId")`: {
				filterParam:   `gt(policyId, "someId")`,
				expectedQuery: bson.D{{"policyId", bson.D{{"$gt", "someId"}}}},
			},
			`given gt(cores, 4)`: {
				filterParam:   `gt(cores, 4)`,
				expectedQuery: bson.D{{"cores", bson.D{{"$gt", 4}}}},
			},
			`given ge(policyId, "someId")`: {
				filterParam:   `ge(policyId, "someId")`,
				expectedQuery: bson.D{{"policyId", bson.D{{"$gte", "someId"}}}},
			},
			`given ge(cores, 4)`: {
				filterParam:   `ge(cores, 4)`,
				expectedQuery: bson.D{{"cores", bson.D{{"$gte", 4}}}},
			},
			`given lt(policyId, "someId")`: {
				filterParam:   `lt(policyId, "someId")`,
				expectedQuery: bson.D{{"policyId", bson.D{{"$lt", "someId"}}}},
			},
			`given lt(cores, 4)`: {
				filterParam:   `lt(cores, 4)`,
				expectedQuery: bson.D{{"cores", bson.D{{"$lt", 4}}}},
			},
			`given le(policyId, "someId")`: {
				filterParam:   `le(policyId, "someId")`,
				expectedQuery: bson.D{{"policyId", bson.D{{"$lte", "someId"}}}},
			},
			`given le(cores, 4)`: {
				filterParam:   `le(cores, 4)`,
				expectedQuery: bson.D{{"cores", bson.D{{"$lte", 4}}}},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				parser := queryparser.New()
				queryBuilder := New(Props{literalComparisonFields})
				tree, _ := parser.Run(tc.filterParam)

				query := queryBuilder.Run(tree)

				assert.Equal(t, tc.expectedQuery, query)
			})
		}
	})
	t.Run("Logical Functions", func(t *testing.T) {
		tests := map[string]test{
			`given and(eq(policyId, "someId"), eq(cores, 4), or(eq(title, "sith"), eq(email, "revan@jedi.com")), "jedi", 2, true )`: {
				filterParam: `and(eq(policyId, "someId"), eq(cores, 4), or(eq(title, "sith"), eq(email, "revan@jedi.com")), "jedi", 2, true )`,
				expectedQuery: bson.D{
					{"$and",
						bson.A{
							bson.D{{"policyId", "someId"}},
							bson.D{{"cores", 4}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", "sith"}},
									bson.D{{"email", "revan@jedi.com"}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", "jedi"}},
									bson.D{{"email", "jedi"}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", 2}},
									bson.D{{"email", 2}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", true}},
									bson.D{{"email", true}},
								},
							}},
						},
					},
				},
			},
			`given or(eq(policyId, "someId"), eq(cores, 4), and(eq(title, "sith"), eq(email, "revan@jedi.com")), "jedi", 2, true )`: {
				filterParam: `or(eq(policyId, "someId"), eq(cores, 4), and(eq(title, "sith"), eq(email, "revan@jedi.com")), "jedi", 2, true )`,
				expectedQuery: bson.D{
					{"$or",
						bson.A{
							bson.D{{"policyId", "someId"}},
							bson.D{{"cores", 4}},
							bson.D{{"$and",
								bson.A{
									bson.D{{"title", "sith"}},
									bson.D{{"email", "revan@jedi.com"}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", "jedi"}},
									bson.D{{"email", "jedi"}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", 2}},
									bson.D{{"email", 2}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", true}},
									bson.D{{"email", true}},
								},
							}},
						},
					},
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				parser := queryparser.New()
				queryBuilder := New(Props{literalComparisonFields})
				tree, _ := parser.Run(tc.filterParam)

				query := queryBuilder.Run(tree)

				assert.Equal(t, tc.expectedQuery, query)
			})
		}
	})
	t.Run("Literals", func(t *testing.T) {
		tests := map[string]test{
			`given "sith"`: {
				filterParam: `"sith"`,
				expectedQuery: bson.D{
					{"$or",
						bson.A{
							bson.D{{"title", "sith"}},
							bson.D{{"email", "sith"}},
						},
					},
				},
			},
			`given 10`: {
				filterParam: `10`,
				expectedQuery: bson.D{
					{"$or",
						bson.A{
							bson.D{{"title", 10}},
							bson.D{{"email", 10}},
						},
					},
				},
			},
			`given true`: {
				filterParam: `true`,
				expectedQuery: bson.D{
					{"$or",
						bson.A{
							bson.D{{"title", true}},
							bson.D{{"email", true}},
						},
					},
				},
			},
			`given and("sith", 2)`: {
				filterParam: `and("sith", 2)`,
				expectedQuery: bson.D{
					{"$and",
						bson.A{
							bson.D{
								{"$or",
									bson.A{
										bson.D{{"title", "sith"}},
										bson.D{{"email", "sith"}},
									},
								},
							},
							bson.D{
								{"$or",
									bson.A{
										bson.D{{"title", 2}},
										bson.D{{"email", 2}},
									},
								},
							},
						},
					},
				},
			},
			`given or("sith", true)`: {
				filterParam: `or("sith", true)`,
				expectedQuery: bson.D{
					{"$or",
						bson.A{
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", "sith"}},
									bson.D{{"email", "sith"}},
								},
							}},
							bson.D{{"$or",
								bson.A{
									bson.D{{"title", true}},
									bson.D{{"email", true}},
								},
							}},
						},
					},
				},
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				parser := queryparser.New()
				queryBuilder := New(Props{literalComparisonFields})
				tree, _ := parser.Run(tc.filterParam)

				query := queryBuilder.Run(tree)

				assert.Equal(t, tc.expectedQuery, query)
			})
		}
	})
}
