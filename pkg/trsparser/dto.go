package trsparser

import "encoding/json"

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../docs/trs-parser-api.yaml

func NewConstantMonomial(constant int) Monomial {
	monomial := Monomial{}
	monomial.FromConstantMonomial(ConstantMonomial{
		Constant: constant,
	})
	return monomial
}

func NewProductMonomial(factors []Factor) Monomial {
	monomial := Monomial{}
	monomial.FromProductMonomial(ProductMonomial{
		Factors: factors,
	})
	return monomial
}

// TODO: remove
func NewSubexpressionArg(sexpr Subexpression) (interface{}, error) {
	b, err := json.Marshal(sexpr)
	return b, err
}
