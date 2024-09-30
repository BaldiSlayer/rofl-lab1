package trsparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorOnEmptyInput(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse("")

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "система должна содержать хотя бы одно правило переписывания и его интерпретацию", parseError.LlmMessage)
}

func TestParsesBasicTrs(t *testing.T) {
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
				Name:       "a",
			},
		},
	}
	inter := Interpretation{
		Args:      []string{},
		Constants: []int{},
		Monomials: []Monomial{{
			Coefficient: nil,
			Power:       nil,
			Variable:    "a",
		}},
		Name: "f",
	}

	trs, err := Parser{}.Parse(
		`variables = a
f = a
-----
f = a
`,
	)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: []Interpretation{inter},
		Rules:           []Rule{rule},
		Variables:       []string{"a"},
	}, *trs)
}
