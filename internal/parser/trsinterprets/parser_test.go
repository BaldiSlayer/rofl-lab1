package trsinterprets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestSingleConstInterpretation(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
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
			Name:      "f",
			Args:      []string{},
			Monomials: []Monomial{NewConstantMonomial(5)},
		},
	}, interpretations)
}

func TestMultipleConstInterpretations(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
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
			Name:      "f",
			Args:      []string{},
			Monomials: []Monomial{NewConstantMonomial(5)},
		},
		{
			Name:      "g",
			Args:      []string{},
			Monomials: []Monomial{NewConstantMonomial(100)},
		},
	}, interpretations)
}

func TestNoInterpretations(t *testing.T) {
	input := ToInputChannel([]models.Lexem{})

	_, err := NewParser(input, map[string]int{}).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "система должна содержать хотя бы одну интерпретацию", parseError.LlmMessage)
}

func TestNoConstructorName(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexEOL, Str: "\n"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
	})

	_, err := NewParser(input, map[string]int{"f": 0, "g": 0}).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t,
		`неверно задана интерпретация: ожидалось название конструктора, получено "=" (строка 0, символ 0)`,
		parseError.LlmMessage)
}

func TestSingleInterpretation(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
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
			Name: "f",
			Args: []string{"x"},
			Monomials: []Monomial{
				NewProductMonomial([]Factor{{
					Variable:    "x",
					Coefficient: 1,
					Power:       1,
				}}),
				NewConstantMonomial(5),
			},
		},
	}, interpretations)
}

func TestMultipleInterpretations(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
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
	constructorArity := map[string]int{"f": 1, "g": 2}

	interpretations, err := NewParser(input, constructorArity).Parse()

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			Name: "f",
			Args: []string{"x"},
			Monomials: []Monomial{
				NewProductMonomial([]Factor{{
					Variable:    "x",
					Coefficient: 1,
					Power:       1,
				}}),
			},
		},
		{
			Name: "g",
			Args: []string{"x", "y"},
			Monomials: []Monomial{
				NewProductMonomial([]Factor{{
					Variable:    "y",
					Coefficient: 1,
					Power:       5,
				}}),
				NewProductMonomial([]Factor{{
					Variable:    "x",
					Coefficient: 13,
					Power:       1,
				}}),
			},
		},
	}, interpretations)
}

func TestMissingStarSign(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
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

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, `неверно задана интерпретация конструктора f: `+
		`при разборе монома в формате [опциональный коэффициент *] переменная [опциональная степень], `+
		`ожидался знак "*" после коэффициента 5, получено "x" `+
		`(строка 0, символ 0)`, parseError.LlmMessage)
}

func TestUndefinedVariable(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
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

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверно задана интерпретация конструктора f: "+
		"аргумент z не объявлен в левой части выражения, но использован в правой", parseError.LlmMessage)
}

func TestInterpretationArityMismatch(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 2}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "неверная арность конструктора f: ожидалась арность 2, получена арность 1", parseError.LlmMessage)
}

func TestExcessInterpretation(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "конструктор f отсутствует в правилах trs", parseError.LlmMessage)
}

func TestDuplicateInterpretation(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "13"},
		{LexemType: models.LexMUL, Str: "*"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 1}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "интерпретация конструктора f задана повторно, "+
		"хотя каждый конструктор должен иметь только одну интерпретацию", parseError.LlmMessage)
}

func TestDuplicateArgument(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexCOMMA, Str: ","},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 2}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "в интерпретации конструктора f повторно объявлена переменная x, "+
		"хотя каждая переменная должна быть объявлена один раз", parseError.LlmMessage)
}

func TestNoSufficientInterpretation(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 1, "g": 2}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "не хватает интерпретации для конструктора g", parseError.LlmMessage)
}

func TestUnusedArgument(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 1}

	_, err := NewParser(input, constructorArity).Parse()

	require.NoError(t, err)
}

func TestMultipleVariablesInMonomial(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexCOMMA, Str: ","},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexLCB, Str: "{"},
		{LexemType: models.LexNUM, Str: "2"},
		{LexemType: models.LexRCB, Str: "}"},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexMUL, Str: "*"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexADD, Str: "+"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexLCB, Str: "{"},
		{LexemType: models.LexNUM, Str: "2"},
		{LexemType: models.LexRCB, Str: "}"},
		{LexemType: models.LexNUM, Str: "5"},
		{LexemType: models.LexMUL, Str: "*"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexLETTER, Str: "y"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 2}

	interpretations, err := NewParser(input, constructorArity).Parse()

	assert.NoError(t, err)
	assert.Equal(t, []Interpretation{
		{
			Name: "f",
			Args: []string{"x", "y"},
			Monomials: []Monomial{
				NewProductMonomial([]Factor{
					{
						Variable:    "x",
						Coefficient: 1,
						Power:       1,
					},
					{
						Variable:    "x",
						Coefficient: 1,
						Power:       2,
					},
					{
						Variable:    "x",
						Coefficient: 5,
						Power:       1,
					},
					{
						Variable:    "y",
						Coefficient: 1,
						Power:       1,
					},
				}),
				NewProductMonomial([]Factor{
					{
						Variable:    "x",
						Coefficient: 1,
						Power:       2,
					},
					{
						Variable:    "x",
						Coefficient: 5,
						Power:       1,
					},
					{
						Variable:    "y",
						Coefficient: 1,
						Power:       1,
					},
				}),
			},
		},
	}, interpretations)
}

func TestIllFormedMonomial(t *testing.T) {
	input := ToInputChannel([]models.Lexem{
		{LexemType: models.LexLETTER, Str: "f"},
		{LexemType: models.LexLB, Str: "("},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexRB, Str: ")"},
		{LexemType: models.LexEQ, Str: "="},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexMUL, Str: "*"},
		{LexemType: models.LexLETTER, Str: "x"},
		{LexemType: models.LexEOL, Str: "\n"},
	})
	constructorArity := map[string]int{"f": 1}

	_, err := NewParser(input, constructorArity).Parse()

	var parseError *models.ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, `неверно задана интерпретация конструктора f: `+
		`при разборе монома в формате [опциональный коэффициент *] переменная [опциональная степень], `+
		`ожидалось название переменной или коэффициент, получено "*" `+
		`(строка 0, символ 0)`, parseError.LlmMessage)
}
