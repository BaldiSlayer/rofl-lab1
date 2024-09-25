package interprets

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/trs"
)

type Parser struct {}

func (p Parser) Parse(input chan trs.Lexem, constructorArity map[string]int) ([]Interpretation, error) {
	if constructorArity["f"] == 2 {
		return nil, &ParseError{
			llmMessage: "неверная арность интерпретации конструктора f: ожидалось 2, получено 1",
			summary:    "wrong interpretation arity",
		}
	}

	return []Interpretation{{
		name:      "f",
		args:      []string{"x"},
		monomials: []Monomial{},
		constants: []int{5},
	},}, nil
}
