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
	expectedMessage, expectedLlmMessage string) (models.Lexem, *ParseError) {

	got := p.stream.next()
	if got.Type() != expectedType {
		return models.Lexem{}, &ParseError{
			llmMessage: fmt.Sprintf("%s, получено %v", expectedLlmMessage, got.String()),
			message:    fmt.Sprintf("expected %s, got %v", expectedMessage, got.String()),
		}
	}
	return got, nil
}

func (p *Parser) peek() models.LexemType {
	return p.stream.peek().Type()
}

func (p *Parser) interpretations() ([]Interpretation, *ParseError) {
	res := []Interpretation{}
	for {
		if p.peek() == models.LexEOF && len(res) == 0 {
			return nil, &ParseError{
				llmMessage: "система должна содержать хотя бы одну интерпретацию",
				message:    "at least one interpretation expected",
			}
		}
		if p.peek() == models.LexEOF {
			return res, nil
		}

		interpret, err := p.interpretation()
		if err != nil {
			return nil, err
		}

		_, err = p.accept(models.LexEOL, "EOL", "ожидался перенос строки после определения интерпретации")
		if err != nil {
			return nil, err.wrap(&ParseError{
				llmMessage: fmt.Sprintf("неверно задана интерпретация %s", interpret.Name),
				message:    "ill-formed interpretation",
			})
		}

		res = append(res, interpret)
	}
}

func (p *Parser) interpretation() (Interpretation, *ParseError) {
	constructor, err := p.accept(
		models.LexLETTER,
		"constructor name",
		"ожидалось название конструктора",
	)
	if err != nil {
		return Interpretation{}, err.wrap(&ParseError{
			llmMessage: "неверно задана интерпретация",
			message:    "wrong interpretation definition",
		})
	}

	interpret, err := p.interpretationBody(constructor.String())
	if err != nil {
		return Interpretation{}, err.wrap(&ParseError{
			llmMessage: fmt.Sprintf("неверно задана интерпретация конструктора %s", constructor.String()),
			message:    "wrong interpretation definition",
		})
	}

	return interpret, nil
}

func (p *Parser) interpretationBody(name string) (Interpretation, *ParseError) {
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
	return Interpretation{}, &ParseError{
		llmMessage: fmt.Sprintf("ожидалось = или ( после названия конструктора, получено %s", got.String()),
		message:    fmt.Sprintf("expected = or (, got %s", got.String()),
	}
}

func (p *Parser) constInterpretation() (int, *ParseError) {
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

func (p *Parser) toInt(lexem models.Lexem) (int, *ParseError) {
	if lexem.Type() != models.LexNUM {
		return 0, &ParseError{
			llmMessage: "ожидалось натуральное число",
			message:    "number expected",
		}
	}
	num, err := strconv.Atoi(lexem.String())
	if err != nil {
		return 0, &ParseError{
			llmMessage: "ошибка в лексере: невозможно сконвертировать лексему числа в число",
			message:    fmt.Sprintf("can't convert number lexem: |%s|", lexem.String()),
		}
	}
	return num, nil
}

func (p *Parser) funcInterpretation(name string) (Interpretation, *ParseError) {
	p.stream.next()

	args, err := p.letters()
	if err != nil {
		return Interpretation{}, err
	}

	p.accept(models.LexRB, ")", "ожидалось закрытие скобки после объявления переменных через запятую")
	p.accept(models.LexEQ, "=", "ожидался знак равенства")

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

func (p *Parser) letters() ([]string, *ParseError) {
	lexem, err := p.accept(models.LexLETTER, "letter", "ожидалась буква - название переменной")
	if err != nil {
		return nil, err
	}

	variables := []string{}
	variables = append(variables, lexem.String())

	for p.peek() == models.LexCOMMA {
		p.stream.next()

		lexem, err := p.accept(models.LexLETTER, "letter", "ожидалась буква - название переменной")
		if err != nil {
			return nil, err
		}

		variables = append(variables, lexem.String())
	}

	return variables, nil
}

func (p *Parser) funcInterpretationBody() ([]Monomial, *ParseError) {
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

func (p *Parser) monomial() (Monomial, *ParseError) {
	monomial, err := p.factorOrConstant()
	if err != nil {
		return Monomial{}, err
	}

	for p.peek() != models.LexADD && p.peek() != models.LexEOL {
		factor, err := p.factor()
		if err != nil {
			return Monomial{}, err
		}

		*monomial.Factors = append(*monomial.Factors, factor)
	}

	return monomial, nil
}

func (p *Parser) factorOrConstant() (Monomial, *ParseError) {
	coefficient := 1

	if p.peek() == models.LexNUM {
		num, err := p.number()
		if err != nil {
			return Monomial{}, err
		}

		if p.peek() == models.LexEOL || p.peek() == models.LexADD {
			return NewConstantMonomial(num), nil
		}

		err = p.starSign(num)
		if err != nil {
			return Monomial{}, err
		}

		coefficient = num
	}

	name, err := p.variable()
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

func (p *Parser) factor() (Factor, *ParseError) {
	coefficient := 1

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
	}

	name, err := p.variable()
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

func (p *Parser) number() (int, *ParseError) {
	numLexem, err := p.accept(models.LexNUM, "expected number in monomial", "ожидалось число в мономе")
	if err != nil {
		return 0, err
	}

	return p.toInt(numLexem)
}

func (p *Parser) starSign(coefficient int) *ParseError {
	_, err := p.accept(
		models.LexMUL,
		"star sign",
		fmt.Sprintf("ожидался знак * после коэффициента %d в определении монома", coefficient),
	)
	return err
}

func (p *Parser) variable() (string, *ParseError) {
	varLexem, err := p.accept(
		models.LexLETTER,
		"variable name",
		"в определении монома ожидалось название переменной или значение коэффициента",
	)
	if err != nil {
		return "", err
	}

	return varLexem.String(), nil
}

func (p *Parser) power() (int, *ParseError) {
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
