package parser

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../docs/trs-parser-api.yaml

func NewConstInterpretation(name string, value int) Interpretation {
	res := Interpretation{}
	res.FromConstInterpretation(ConstInterpretation{
		Name:  name,
		Value: value,
	})
	return res
}
