package trsparser

type Parser struct{}

func (p Parser) Parse(input string) (*Trs, error) {
	if input == "" {
		return nil, &ParseError{
			LlmMessage: "система должна содержать хотя бы одно правило переписывания и его интерпретацию",
			Summary:    "empty input",
		}
	}

	inter := Interpretation{
		Args:      []string{},
		Constants: []int{5},
		Monomials: []Monomial{},
		Name:      "f",
	}
	rule := Rule{
		Lhs: Subexpression{
			Args: nil,
			Letter: Letter{
				IsVariable: false,
				Name:       "f",
			},
		},
		Rhs: Subexpression{
			Args: nil,
			Letter: Letter{
				IsVariable: true,
				Name:       "x",
			},
		},
	}

	return &Trs{
		Interpretations: []Interpretation{inter},
		Rules:           []Rule{rule},
		Variables:       []string{"x"},
	}, nil
}
