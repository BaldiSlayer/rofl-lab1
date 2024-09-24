package parser

import "fmt"

type Parser struct {}

func (p Parser) Parse(input string) (*Trs, error) {
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	inter := Interpretation{
		Args:      []string{},
		Constants: []int{5},
		Monomials: []Monomial{},
		Name:      "f",
	}
	rule := Rule{
		Lhs: Subexpression{
			Args:   nil,
			Letter: Letter{
				IsVariable: false,
				Name:       "f",
			},
		},
		Rhs: Subexpression{
			Args:   nil,
			Letter: Letter{
				IsVariable: true,
				Name:       "a",
			},
		},
	}

	return &Trs{
		Interpretations: []Interpretation{inter},
		Rules:           []Rule{rule},
		Variables:       []string{"a"},
	}, nil
}
