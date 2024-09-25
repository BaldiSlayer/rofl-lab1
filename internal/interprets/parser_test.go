package interprets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type lexem struct {
	t int
	v string
}

func (l lexem) String() string {
	return l.v
}

func (l lexem) Type() int {
	return l.t
}

func TestInterpretationArityMismatch(t *testing.T) {
	// f(x) = 5
	lexems := []lexem{
		{lex_LETTER, "f"},
		{lex_LB, "("},
		{lex_LETTER, "x"},
		{lex_RB, ")"},
		{lex_EQ, "="},
		{lex_NUM, "5"},
	}
	constructorArity := map[string]int{"f": 2}
	input := make(chan Lexem, 100)
	for _, el := range lexems {
		input <- el
	}

	_, err := Parser{}.Parse(input, constructorArity)

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверная арность интерпретации конструктора f: ожидалось 2, получено 1", parseError.LlmMessage())
}

func TestSingleInterpretation(t *testing.T) {
	// f(x) = 5
	lexems := []lexem{
		{lex_LETTER, "f"},
		{lex_LB, "("},
		{lex_LETTER, "x"},
		{lex_RB, ")"},
		{lex_EQ, "="},
		{lex_NUM, "5"},
	}
	constructorArity := map[string]int{"f": 1}
	input := make(chan Lexem, 100)
	for _, el := range lexems {
		input <- el
	}

	interpretations, err := Parser{}.Parse(input, constructorArity)

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			name:      "f",
			args:      []string{"x"},
			monomials: []Monomial{},
			constants: []int{5},
		},
	}, interpretations)
}
