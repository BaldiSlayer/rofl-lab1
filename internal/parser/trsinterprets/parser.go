package trsinterprets

import (
	"fmt"
	"strconv"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
)

type Parser struct {
	stream           stream
	ConstructorArity map[string]int
}

func NewParser(input chan models.Lexem, constructorArity map[string]int) *Parser {
	return &Parser{
		stream: stream{
			channel: input,
			ok:      false,
			val:     models.Lexem{},
		},
		ConstructorArity: constructorArity,
	}
}

func (p *Parser) Parse() ([]Interpretation, error) {
	i, err := p.interpretations()
	if err != nil {
		return nil, err
	}

	err = checkInterpretations(i, p.ConstructorArity)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (p *Parser) accept(expectedType models.LexemType,
	expectedMessage, expectedLlmMessage string) (models.Lexem, error) {

	got := p.stream.next()
	if got.Type() != expectedType {
		return models.Lexem{}, &models.ParseError{
			LlmMessage: fmt.Sprintf(
				`%s, получено "%s" (строка %d, символ %d)`,
				expectedLlmMessage,
				got.String(),
				got.Line,
				got.Index,
			),
			Message: fmt.Sprintf("expected %s, got %v", expectedMessage, got.String()),
		}
	}
	return got, nil
}

func (p *Parser) peek() models.LexemType {
	return p.stream.peek().Type()
}

func (p *Parser) interpretations() ([]Interpretation, error) {
	res := []Interpretation{}
	for p.peek() != models.LexEOF {
		interpret, err := p.interpretation()
		if err != nil {
			return nil, err
		}

		res = append(res, interpret)

		if p.peek() != models.LexEOF {
			_, err = p.accept(models.LexEOL, "EOL", "ожидался перенос строки после определения интерпретации")
			if err != nil {
				return nil, models.Wrap("ill-formed interpretation", fmt.Sprintf("неверно задана интерпретация %s", interpret.Name), err)
			}
		}
	}

	if len(res) == 0 {
		return nil, &models.ParseError{
			LlmMessage: "система должна содержать хотя бы одну интерпретацию",
			Message:    "at least one interpretation expected",
		}
	}
	return res, nil
}

func (p *Parser) interpretation() (Interpretation, error) {
	constructor, err := p.accept(
		models.LexLETTER,
		"constructor name",
		"ожидалось название конструктора",
	)
	if err != nil {
		return Interpretation{}, models.Wrap("wrong interpretation definition", "неверно задана интерпретация", err)
	}

	interpret, err := p.interpretationBody(constructor.String())
	if err != nil {
		return Interpretation{}, models.Wrap(
			"wrong interpretation definition",
			fmt.Sprintf("неверно задана интерпретация конструктора %s", constructor.String()),
			err,
		)
	}

	return interpret, nil
}

func (p *Parser) interpretationBody(name string) (Interpretation, error) {
	switch p.peek() {
	case models.LexEQ:
		value, err := p.constInterpretation()
		return Interpretation{
			Name: name,
			Args: []string{},
			Monomials: []Monomial{{
				Constant: &value,
				Factors:  nil,
			}},
		}, err
	case models.LexLB:
		return p.funcInterpretation(name)
	}

	got := p.stream.next()
	return Interpretation{}, &models.ParseError{
		LlmMessage: fmt.Sprintf("ожидалось = или ( после названия конструктора, получено %s", got.String()),
		Message:    fmt.Sprintf("expected = or (, got %s", got.String()),
	}
}

func (p *Parser) constInterpretation() (int, error) {
	p.stream.next()
	lexem, err := p.accept(models.LexNUM,
		"expected number in const interpretation",
		"ожидалось натуральное число после знака = в интерпретации константы",
	)
	if err != nil {
		return 0, err
	}
	num, err := p.toInt(lexem)
	return num, err
}

func (p *Parser) toInt(lexem models.Lexem) (int, error) {
	if lexem.Type() != models.LexNUM {
		return 0, &models.ParseError{
			LlmMessage: "ожидалось натуральное число",
			Message:    "number expected",
		}
	}
	num, err := strconv.Atoi(lexem.String())
	if err != nil {
		return 0, &models.ParseError{
			LlmMessage: "ошибка в лексере: невозможно сконвертировать лексему числа в число",
			Message:    fmt.Sprintf("can't convert number lexem: |%s|", lexem.String()),
		}
	}
	return num, nil
}

func (p *Parser) funcInterpretation(name string) (Interpretation, error) {
	p.stream.next()

	args, err := p.letters()
	if err != nil {
		return Interpretation{}, err
	}

	_, err = p.accept(models.LexRB, ")", "при разборе определения переменных через запятую, "+
		"ожидалась запятая или закрывающая скобка после перечисления переменных")
	if err != nil {
		return Interpretation{}, err
	}
	_, err = p.accept(models.LexEQ, "=", "ожидался знак равенства после объявления переменных")
	if err != nil {
		return Interpretation{}, err
	}

	monomials, err := p.funcInterpretationBody()
	if err != nil {
		return Interpretation{}, err
	}

	err = checkMonomials(monomials, args)
	if err != nil {
		return Interpretation{}, err
	}

	return Interpretation{
		Name:      name,
		Args:      args,
		Monomials: monomials,
	}, nil
}

func (p *Parser) letters() ([]string, error) {
	lexem, err := p.accept(
		models.LexLETTER,
		"letter",
		"в объявлении переменных ожидалась хотя бы одна буква - название переменной",
	)
	if err != nil {
		return nil, err
	}

	variables := []string{}
	variables = append(variables, lexem.String())

	for p.peek() == models.LexCOMMA {
		p.stream.next()

		lexem, err := p.accept(
			models.LexLETTER,
			"letter",
			"в объявлении переменных после запятой ожидалась буква - название переменной",
		)
		if err != nil {
			return nil, err
		}

		variables = append(variables, lexem.String())
	}

	return variables, nil
}

func (p *Parser) funcInterpretationBody() ([]Monomial, error) {
	monomials := []Monomial{}

	monomial, err := p.monomial()
	if err != nil {
		return nil, err
	}

	monomials = append(monomials, monomial)

	for p.peek() == models.LexADD {
		p.stream.next()

		monomial, err = p.monomial()
		if err != nil {
			return nil, err
		}

		monomials = append(monomials, monomial)
	}

	return monomials, nil
}

func (p *Parser) monomial() (Monomial, error) {
	monomial, err := p.factorOrConstant()
	if err != nil {
		return Monomial{}, err
	}

	for p.peek() != models.LexADD && p.peek() != models.LexEOL && p.peek() != models.LexEOF {
		factor, err := p.factor()
		if err != nil {
			return Monomial{}, err
		}

		*monomial.Factors = append(*monomial.Factors, factor)
	}

	return monomial, nil
}

func (p *Parser) factorOrConstant() (Monomial, error) {
	coefficient := 1
	coefficientRead := false

	if p.peek() == models.LexNUM {
		num, err := p.number()
		if err != nil {
			return Monomial{}, err
		}

		if p.peek() == models.LexEOL || p.peek() == models.LexADD || p.peek() == models.LexEOF {
			return NewConstantMonomial(num), nil
		}

		err = p.starSign(num)
		if err != nil {
			return Monomial{}, err
		}

		coefficient = num
		coefficientRead = true
	}

	name, err := p.variable(!coefficientRead)
	if err != nil {
		return Monomial{}, err
	}

	power, err := p.power()
	if err != nil {
		return Monomial{}, err
	}

	return NewProductMonomial([]Factor{{
		Variable:    name,
		Coefficient: coefficient,
		Power:       power,
	}}), nil
}

func (p *Parser) factor() (Factor, error) {
	coefficient := 1
	coefficientRead := false

	if p.peek() == models.LexNUM {
		num, err := p.number()
		if err != nil {
			return Factor{}, err
		}

		err = p.starSign(num)
		if err != nil {
			return Factor{}, err
		}

		coefficient = num
		coefficientRead = true
	}

	name, err := p.variable(!coefficientRead)
	if err != nil {
		return Factor{}, err
	}

	power, err := p.power()
	if err != nil {
		return Factor{}, err
	}

	return Factor{
		Variable:    name,
		Coefficient: coefficient,
		Power:       power,
	}, nil
}

func (p *Parser) number() (int, error) {
	numLexem, err := p.accept(models.LexNUM, "expected number in monomial", "ожидалось число в мономе")
	if err != nil {
		return 0, err
	}

	return p.toInt(numLexem)
}

func (p *Parser) starSign(coefficient int) error {
	_, err := p.accept(
		models.LexMUL,
		"star sign",
		fmt.Sprintf("при разборе монома в формате [опциональный коэффициент *] переменная [опциональная степень], "+
			"ожидался знак \"*\" после коэффициента %d", coefficient),
	)
	return err
}

func (p *Parser) variable(expectCoefficient bool) (string, error) {
	var expectCoeffText string
	if expectCoefficient {
		expectCoeffText = " или коэффициент"
	}

	varLexem, err := p.accept(
		models.LexLETTER,
		"variable name",
		fmt.Sprintf("при разборе монома в формате [опциональный коэффициент *] переменная [опциональная степень], "+
			"ожидалось название переменной%s",
			expectCoeffText,
		),
	)
	if err != nil {
		return "", err
	}

	return varLexem.String(), nil
}

func (p *Parser) power() (int, error) {
	if p.peek() != models.LexLCB {
		return 1, nil
	}

	p.stream.next()
	numLexem, err := p.accept(
		models.LexNUM,
		"number expected in power definition",
		"после { ожидалось значение степени - натуральное число",
	)
	if err != nil {
		return 0, err
	}
	num, err := p.toInt(numLexem)
	if err != nil {
		return 0, err
	}

	_, err = p.accept(models.LexRCB, "}", "ожидалось закрытие фигурных скобок }")
	if err != nil {
		return 0, err
	}

	return num, err
}

func ToInputChannel(lexems []models.Lexem) chan models.Lexem {
	channel := make(chan models.Lexem, 100)
	for _, el := range lexems {
		channel <- el
	}
	close(channel)
	return channel
}
