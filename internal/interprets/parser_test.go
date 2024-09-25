package interprets

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BaldiSlayer/rofl-lab1/internal/trs"
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

func toInputChannel(lexems []lexem) chan trs.Lexem {
	channel := make(chan trs.Lexem, 100)
	for _, el := range lexems {
		channel <- el
	}
	return channel
}

func TestInterpretationArityMismatch(t *testing.T) {
	// f(x) = 5
	input := toInputChannel([]lexem{
		{trs.LexLETTER, "f"},
		{trs.LexLB, "("},
		{trs.LexLETTER, "x"},
		{trs.LexRB, ")"},
		{trs.LexEQ, "="},
		{trs.LexNUM, "5"},
	})
	constructorArity := map[string]int{"f": 2}

	_, err := Parser{}.Parse(input, constructorArity)

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверная арность интерпретации конструктора f: ожидалось 2, получено 1", parseError.LlmMessage())
}

func TestSingleInterpretation(t *testing.T) {
	// f(x) = 5
	input := toInputChannel([]lexem{
		{trs.LexLETTER, "f"},
		{trs.LexLB, "("},
		{trs.LexLETTER, "x"},
		{trs.LexRB, ")"},
		{trs.LexEQ, "="},
		{trs.LexNUM, "5"},
	})
	constructorArity := map[string]int{"f": 1}

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
