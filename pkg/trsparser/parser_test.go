package trsparser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorOnEmptyInput(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse("")

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		"система должна содержать хотя бы одно правило переписывания и его интерпретацию",
		parseError.LlmMessage,
	)
}

func TestParsesBasicTrs(t *testing.T) {
	expectedRule := Rule{
		Lhs: Subexpression{
			Args: &[]interface{}{},
			Letter: Letter{
				IsVariable: false,
				Name:       "f",
			},
		},
		Rhs: Subexpression{
			Args: &[]interface{}{},
			Letter: Letter{
				IsVariable: true,
				Name:       "x",
			},
		},
	}
	expectedInterpretation := Interpretation{
		Args:      []string{},
		Monomials: []Monomial{NewConstantMonomial(5)},
		Name:      "f",
	}

	trs, err := Parser{}.Parse(
		`variables = x
f = x
-----
f = 5
`,
	)

	b, _ := json.Marshal(*trs)
	os.Stdout.Write(b)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: []Interpretation{expectedInterpretation},
		Rules:           []Rule{expectedRule},
		Variables:       []string{"x"},
	}, *trs)
}

func TestParsesComplexTrs(t *testing.T) {
	expectedRule := Rule{
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
	expectedInterpretations := []Interpretation{
		{
			Args: []string{"x", "y"},
			Monomials: []Monomial{
				NewProductMonomial(
					[]Factor{{
						Coefficient: newInt(5),
						Power:       newInt(2),
						Variable:    "x",
					}},
				),
			},
			Name: "f",
		},
		{
			Args:      []string{},
			Monomials: []Monomial{},
			Name:      "f",
		},
	}

	trs, err := Parser{}.Parse(
		`variables = x, y
f(x, g(y)) = g(f(x, h))
f(y, x) = g(y)
-----
f(x, y) = 5*x{2} + 10 + y{120}
g(x) = xx{2}5*x
h = 123
`,
	)

	// b, _ := json.Marshal(*trs)
	// os.Stdout.Write(b)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: expectedInterpretations,
		Rules:           []Rule{expectedRule},
		Variables:       []string{"x"},
	}, *trs)
}

func newInt(v int) *int {
	return &v
}
