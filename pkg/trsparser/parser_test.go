package trsparser

import (
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

func TestErrorOnJustEOL(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(`

`)
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`в начале TRS ожидалось перечисление переменных формата "variables = x,y,z"`,
		parseError.LlmMessage,
	)
}

func TestErrorNoEOL(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
-----
f(x) = 5`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация f: `+
			`ожидался перенос строки после определения интерпретации, получено "EOF" (строка 0, символ 0)`,
		parseError.LlmMessage,
	)
}

func TestErrorNoSeparator(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
f(x) = 5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`в строке 3 TRS  на позиции 8 ожидалось "буква", найдено "5"`,
		parseError.LlmMessage,
	)
}

func TestExcessVariables(t *testing.T) {
	_, err := Parser{}.Parse(
		`variables = x, y
f(x) = x
------
f(x) = 5
`,
	)

	assert.NoError(t, err)
}

func TestMissingEqualsSign(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
------
f(x) 5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`ожидался знак равенства после объявления переменных, получено "5" (строка 4, символ 6)`,
		parseError.LlmMessage,
	)
}

func TestExcessClosingBracket(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
------
f(x)) = 5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`ожидался знак равенства после объявления переменных, получено ")" (строка 4, символ 5)`,
		parseError.LlmMessage,
	)
}

func TestExcessStarSign(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
------
f(x) = 5**x
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`при разборе монома в формате [опциональный коэффициент *] переменная [опциональная степень], `+
			`ожидалось название переменной, получено "*" (строка 4, символ 10)`,
		parseError.LlmMessage,
	)
}

func TestMissingClosingBracket(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
------
f(x = 5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`при разборе определения переменных через запятую, `+
			`ожидалась запятая или закрывающая скобка после перечисления переменных, получено "=" `+
			`(строка 4, символ 5)`,
		parseError.LlmMessage,
	)
}

func TestExcessClosingBracketInRules(t *testing.T) {
	//var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x = x
-----
f(x) = 5
`,
	)
	assert.NoError(t, err)

}

func TestMissingEqualSignAtVariablesBlock(t *testing.T) {
	//var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables x
f(x) = x
-----
f(x) = 5
`,
	)
	assert.NoError(t, err)
}

func TestCoefficientAfterVariable(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
-----
f(x) = x*5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`при разборе монома в формате [опциональный коэффициент *] переменная [опциональная степень], `+
			`ожидалось название переменной или коэффициент, получено "*" (строка 4, символ 9)`,
		parseError.LlmMessage,
	)
}

func TestNoVariablesInInterpretation(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
-----
f() = x*5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`в объявлении переменных ожидалась хотя бы одна буква - название переменной, получено ")" `+
			`(строка 4, символ 3)`,
		parseError.LlmMessage,
	)
}

func TestNestedBracketsInInterpretation(t *testing.T) {
	var parseError *ParseError

	_, err := Parser{}.Parse(
		`variables = x
f(x) = x
-----
f(x, g(y)) = x*5
`,
	)

	assert.ErrorAs(t, err, &parseError)
	assert.Equal(
		t,
		`неверно задана интерпретация конструктора f: `+
			`при разборе определения переменных через запятую, `+
			`ожидалась запятая или закрывающая скобка после перечисления переменных, получено "(" `+
			`(строка 4, символ 7)`,
		parseError.LlmMessage,
	)
}

func TestParsesSimpleTrs(t *testing.T) {
	trs, err := Parser{}.Parse(
		`variables = x
f(x) = x
-----
f(x) = 5
`,
	)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: []string{"f(x)=5"},
		Rules:           []string{"f(x)=x"},
		Variables:       []string{"x"},
	}, *trs)
}

func TestParsesComplexTrs(t *testing.T) {
	trs, err := Parser{}.Parse(
		`variables = x, y

f(x, g(y)) = g(f(x, h))
f(y, x) = g(y)

-----
f(x, y) = 6*x{322} + 10 + y{120}
g(x) = xx{2}5*x
h = 123
`,
	)

	assert.NoError(t, err)
	assert.Equal(t, Trs{
		Interpretations: []string{
			"f(x,y)=6*x**322+10+y**120",
			"g(x)=x*x**2*5*x",
			"h=123",
		},
		Rules: []string{
			"f(x,g(y))=g(f(x,h))",
			"f(y,x)=g(y)",
		},
		Variables: []string{"x", "y"},
	}, *trs)
}

func TestSpecificRules(t *testing.T) {
	/*исходный тест исправленный по ошибкам
		`variables=x
	f(x,g(x,y,y) = f(x,y)
	----------
	f(x,y,z) = x+y
	y = 0
	g(x,y) = x*y`
	*/

	_, err := Parser{}.Parse(
		`variables=x
f(x,g(x,y),y) = f(x,y,y)
----------
f(x,y,z) = x+y
y = 0
g(x,y) = xy
`,
	)
	assert.NoError(t, err)
}
