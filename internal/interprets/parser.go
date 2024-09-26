package interprets

import (
	"fmt"
	"strconv"

	"github.com/BaldiSlayer/rofl-lab1/internal/trs"
)

type Parser struct {
	stream           stream
	ConstructorArity map[string]int
}

func NewParser(input chan trs.Lexem, constructorArity map[string]int) *Parser {
	return &Parser{
		stream: stream{
			channel: input,
			ok:      false,
			val:     trs.Lexem{},
		},
		ConstructorArity: constructorArity,
	}
}

func (p *Parser) Parse() ([]Interpretation, error) {
	i, err := p.interpretations()
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (p *Parser) accept(expectedType trs.LexemType, expectedMessage, expectedLlmMessage string) (trs.Lexem, *ParseError) {
	got := p.stream.next()
	if got.Type() != expectedType {
		return trs.Lexem{}, &ParseError{
			llmMessage: fmt.Sprintf("%s, получено %v", expectedLlmMessage, got.String()),
			message:    fmt.Sprintf("expected %s, got %v", expectedMessage, got.String()),
		}
	}
	return got, nil
}

func (p *Parser) peek() trs.LexemType {
	return p.stream.peek().Type()
}

func (p *Parser) interpretations() ([]Interpretation, *ParseError) {
	res := []Interpretation{}
	for {
		if p.peek() == trs.LexEOF && len(res) == 0 {
			return nil, &ParseError{
				llmMessage: "система должна содержать хотя бы одну интерпретацию",
				message:    "at least one interpretation expected",
			}
		}
		if p.peek() == trs.LexEOF {
			return res, nil
		}

		interpret, err := p.interpretation()
		if err != nil {
			return nil, err
		}

		p.accept(trs.LexEOL, "EOL", "ожидался перенос строки после определения интерпретации")

		res = append(res, interpret)
	}
}

func (p *Parser) interpretation() (Interpretation, *ParseError) {
	constructor, err := p.accept(
		trs.LexLETTER,
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
	case trs.LexEQ:
		value, err := p.constInterpretation()
		return Interpretation{
			name:      name,
			args:      []string{},
			monomials: []Monomial{},
			constants: []int{value},
		}, err
	case trs.LexLB:
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
	lexem, err := p.accept(trs.LexNUM, "number", "ожидалось натуральное число после знака = в интерпретации константы")
	if err != nil {
		return 0, err
	}
	num, err := p.toInt(lexem)
	return num, err
}

func (p *Parser) toInt(lexem trs.Lexem) (int, *ParseError) {
	num, err := strconv.Atoi(lexem.String())
	if err != nil || lexem.Type() != trs.LexNUM {
		return 0, &ParseError{
			llmMessage: "ожидалось натуральное число",
			message:    "number",
		}
	}
	return num, nil
}

func (p *Parser) funcInterpretation(name string) (Interpretation, *ParseError) {
	p.stream.next()

	// TODO: check if name occures in args
	// TODO: check for duplicate args
	args, err := p.letters()
	if err != nil {
		return Interpretation{}, err
	}

	p.accept(trs.LexRB, ")", "ожидалось закрытие скобки после объявления переменных через запятую")
	p.accept(trs.LexEQ, "=", "ожидался знак равенства")

	monomials, constants, err := p.funcInterpretationBody()
	if err != nil {
		return Interpretation{}, err
	}

	return Interpretation{
		name:      name,
		args:      args,
		monomials: monomials,
		constants: constants,
	}, nil
}

func (p *Parser) letters() ([]string, *ParseError) {
	lexem, err := p.accept(trs.LexLETTER, "letter", "ожидалась буква - название переменной")
	if err != nil {
		return nil, err
	}

	variables := []string{}
	variables = append(variables, lexem.String())

	for p.peek() == trs.LexCOMMA {
		p.stream.next()

		lexem, err := p.accept(trs.LexLETTER, "letter", "ожидалась буква - название переменной")
		if err != nil {
			return nil, err
		}

		variables = append(variables, lexem.String())
	}

	return variables, nil
}

func (p *Parser) funcInterpretationBody() ([]Monomial, []int, *ParseError) {
	monomials := []Monomial{}
	constants := []int{}

	appendIfNotNil := func(monomial *Monomial, constant *int) {
		if monomial != nil {
			monomials = append(monomials, *monomial)
		}
		if constant != nil {
			constants = append(constants, *constant)
		}
	}

	monomial, constant, err := p.monomialOrConstant()
	if err != nil {
		return nil, nil, err
	}
	appendIfNotNil(monomial, constant)

	for p.peek() == trs.LexADD {
		p.stream.next()

		monomial, constant, err = p.monomialOrConstant()
		if err != nil {
			return nil, nil, err
		}
		appendIfNotNil(monomial, constant)
	}

	return monomials, constants, nil
}

func (p *Parser) monomialOrConstant() (*Monomial, *int, *ParseError) {
	coefficient := 1

	if p.peek() == trs.LexNUM {
		numLexem := p.stream.next()
		num, err := p.toInt(numLexem)
		if err != nil {
			return nil, nil, err
		}

		if p.peek() != trs.LexMUL {
			// read constant
			return nil, &num, nil
		}

		p.stream.next()
		coefficient = num
	}

	name, err := p.variable()
	if err != nil {
		return nil, nil, err
	}

	power, err := p.power()
	if err != nil {
		return nil, nil, err
	}

	return &Monomial{
		variable:    name,
		coefficient: coefficient,
		power:       power,
	}, nil, nil
}

func (p *Parser) variable() (string, *ParseError) {
	varLexem, err := p.accept(trs.LexLETTER, "variable name", "ожидалось название переменной")
	if err != nil {
		return "", err
	}
	return varLexem.String(), nil
}

func (p *Parser) power() (int, *ParseError) {
	if p.peek() != trs.LexLCB {
		return 1, nil
	}

	p.stream.next()
	numLexem, err := p.accept(trs.LexNUM, "number", "после { ожидалось значение степени - натуральное число")
	if err != nil {
		return 0, err
	}
	num, err := p.toInt(numLexem)
	if err != nil {
		return 0, err
	}

	_, err = p.accept(trs.LexRCB, "}", "ожидалось закрытие фигурных скобок }")
	if err != nil {
		return 0, err
	}

	return num, err
}
