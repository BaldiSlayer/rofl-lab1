package interprets

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
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

func toInputChannel(lexems []models.Lexem) chan models.Lexem {
	channel := make(chan models.Lexem, 100)
	for _, el := range lexems {
		channel <- el
	}
	close(channel)
	return channel
}

func TestSingleConstInterpretation(t *testing.T) {
	// f = 5
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexEOL, Str: "\n"},
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
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexEOL, Str: "\\n"},
		{LexemType: models.LexLETTER, Str: "g"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "100"},
		{LexemType: models.LexEOL, Str: "\n"},
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
			constants: []int{100},
		},
	}, interpretations)
}

func TestNoInterpretations(t *testing.T) {
	input := toInputChannel([]models.Lexem{})

	_, err := NewParser(input, map[string]int{}).Parse()

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "система должна содержать хотя бы одну интерпретацию", parseError.LlmMessage())
}

func TestNoConstructorName(t *testing.T) {
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexEOL, Str: "\n"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
	})

	_, err := NewParser(input, map[string]int{"f": 0, "g": 0}).Parse()

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверно задана интерпретация: ожидалось название конструктора, получено =", parseError.LlmMessage())
}

func TestSingleInterpretation(t *testing.T) {
	// f(x) = 5
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexADD, Str: "+"},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 1}

	interpretations, err := NewParser(input, constructorArity).Parse()

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			name: "f",
			args: []string{"x"},
			monomials: []Monomial{{
				variable:    "x",
				coefficient: 1,
				power:       1,
			}},
			constants: []int{5},
		},
	}, interpretations)
}

func TestMultipleInterpretations(t *testing.T) {
	// f(x) = 5
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
		{LexemType: models.LexLETTER, Str: "g"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexCOMMA, Str: ","},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexLCB, Str: "{"},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexRCB, Str: "}"},
		{LexemType: models.LexADD, Str: "+"},
		{LexemType: models.LexNUM, Str: "13"},
		{LexemType: models.LexMUL, Str: "*"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 1, "g": 1}

	interpretations, err := NewParser(input, constructorArity).Parse()

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			name: "f",
			args: []string{"x"},
			monomials: []Monomial{{
				variable:    "x",
				coefficient: 1,
				power:       1,
			}},
			constants: []int{},
		},
		{
			name: "g",
			args: []string{"x", "y"},
			monomials: []Monomial{
				{
					variable:    "y",
					coefficient: 1,
					power:       5,
				},
				{
					variable:    "x",
					coefficient: 13,
					power:       1,
				},
			},
			constants: []int{},
		},
	}, interpretations)
}

func TestMissingStarSign(t *testing.T) {
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexCOMMA, Str: ","},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexLCB, Str: "{"},
		{LexemType: models.LexNUM, Str: "10"},
		{LexemType: models.LexRCB, Str: "}"},
		{LexemType: models.LexEOL, Str: "\n"},
	})

	_, err := NewParser(input, map[string]int{"f": 2}).Parse()

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверно задана интерпретация конструктора f: "+
		"ожидался знак * после коэффициента 5 в определении монома, получено x", parseError.LlmMessage())
}

func TestUndefinedVariable(t *testing.T) {
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexCOMMA, Str: ","},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexMUL, Str: "*"},
		{LexemType: models.LexLETTER, Str: "z"},
		{LexemType: models.LexLCB, Str: "{"},
		{LexemType: models.LexNUM, Str: "2"},
		{LexemType: models.LexRCB, Str: "}"},
		{LexemType: models.LexEOL, Str: "\n"},
	})

	_, err := NewParser(input, map[string]int{"f": 2}).Parse()

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверно задана интерпретация конструктора f: "+
		"не объявлен аргумент z", parseError.LlmMessage())
}

func TestInterpretationArityMismatch(t *testing.T) {
	t.SkipNow()
	// f(x) = 5
	input := toInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
	})
	constructorArity := map[string]int{"f": 2}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверная арность интерпретации конструктора f: ожидалось 2, получено 1", parseError.LlmMessage())
}
