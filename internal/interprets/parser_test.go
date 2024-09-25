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

func toInputChannel(lexems []trs.Lexem) chan trs.Lexem {
	channel := make(chan trs.Lexem, 100)
	for _, el := range lexems {
		channel <- el
	}
	close(channel)
	return channel
}

func TestSingleConstInterpretation(t *testing.T) {
	// f = 5
	input := toInputChannel([]trs.Lexem{
		{LexemType: trs.LexLETTER, Str: "f"},
		{LexemType: trs.LexEQ, Str: "="},
		{LexemType: trs.LexNUM, Str: "5"},
	})
	constructorArity := map[string]int{"f": 0}

	interpretations, err := NewParser(input, constructorArity).Parse()

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			name:      "f",
			args:      []string{},
			monomials: []Monomial{},
			constants: []int{5},
		},
	}, interpretations)
}

func TestMultipleConstInterpretations(t *testing.T) {
	// f = 5
	input := toInputChannel([]trs.Lexem{
		{LexemType: trs.LexLETTER, Str: "f"},
		{LexemType: trs.LexEQ, Str: "="},
		{LexemType: trs.LexNUM, Str: "5"},
		{LexemType: trs.LexEOL, Str: "\\n"},
		{LexemType: trs.LexLETTER, Str: "g"},
		{LexemType: trs.LexEQ, Str: "="},
		{LexemType: trs.LexNUM, Str: "1"},
	})
	constructorArity := map[string]int{"f": 0, "g": 0}

	interpretations, err := NewParser(input, constructorArity).Parse()

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			name:      "f",
			args:      []string{},
			monomials: []Monomial{},
			constants: []int{5},
		},
		{
			name:      "g",
			args:      []string{},
			monomials: []Monomial{},
			constants: []int{1},
		},
	}, interpretations)
}

func TestInterpretationArityMismatch(t *testing.T) {
	t.SkipNow()
	// f(x) = 5
	input := toInputChannel([]trs.Lexem{
		{LexemType: trs.LexLETTER, Str: "f"},
		{LexemType: trs.LexLB, Str: "("},
		{LexemType: trs.LexLETTER, Str: "x"},
		{LexemType: trs.LexRB, Str: ")"},
		{LexemType: trs.LexEQ, Str: "="},
		{LexemType: trs.LexNUM, Str: "5"},
	})
	constructorArity := map[string]int{"f": 2}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверная арность интерпретации конструктора f: ожидалось 2, получено 1", parseError.LlmMessage())
}

func TestSingleInterpretation(t *testing.T) {
	t.SkipNow()
	// f(x) = 5
	input := toInputChannel([]trs.Lexem{
		{LexemType: trs.LexLETTER, Str: "f"},
		{LexemType: trs.LexLB, Str: "("},
		{LexemType: trs.LexLETTER, Str: "x"},
		{LexemType: trs.LexRB, Str: ")"},
		{LexemType: trs.LexEQ, Str: "="},
		{LexemType: trs.LexNUM, Str: "5"},
	})
	constructorArity := map[string]int{"f": 1}

	interpretations, err := NewParser(input, constructorArity).Parse()

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
