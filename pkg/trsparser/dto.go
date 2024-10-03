package trsparser

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

func NewSubexpression(sexpr Subexpression) interface{} {
	return sexpr
}
