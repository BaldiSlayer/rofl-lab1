package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorOnEmptyInput(t *testing.T) {
	_, err := Parser{}.Parse("")

	require.Error(t, err)
}

func TestParsesBasicTrs(t *testing.T) {
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
	inter := Interpretation{
		Args:      []string{},
		Constants: []int{5},
		Monomials: []Monomial{},
		Name:      "f",
	}

	trs, err := Parser{}.Parse(
`variables = a
f = a
-----
f = 5
`,
	)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: []Interpretation{inter},
		Rules:           []Rule{rule},
		Variables:       []string{"a"},
		}, *trs)
}
