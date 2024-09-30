package interprets

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
				llmMessage: fmt.Sprintf("неверно задана интерпретация %s", interpret.name),
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
			name:      name,
			args:      []string{},
			monomials: []Monomial{},
			constants: []int{value},
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
		"number",
		"ожидалось натуральное число после знака = в интерпретации константы",
	)
	if err != nil {
		return 0, err
	}
	num, err := p.toInt(lexem)
	return num, err
}

func (p *Parser) toInt(lexem models.Lexem) (int, *ParseError) {
	num, err := strconv.Atoi(lexem.String())
	if err != nil || lexem.Type() != models.LexNUM {
		return 0, &ParseError{
			llmMessage: "ожидалось натуральное число",
			message:    "number",
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

	monomials, constants, err := p.funcInterpretationBody()
	if err != nil {
		return Interpretation{}, err
	}

	err = checkMonomials(monomials, args)
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

	for p.peek() == models.LexADD {
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

	if p.peek() == models.LexNUM {
		numLexem := p.stream.next()
		num, err := p.toInt(numLexem)
		if err != nil {
			return nil, nil, err
		}

		if p.peek() == models.LexEOL || p.peek() == models.LexADD {
			return nil, &num, nil
		}

		_, err = p.accept(
			models.LexMUL,
			"star sign",
			fmt.Sprintf("ожидался знак * после коэффициента %d в определении монома", num),
		)
		if err != nil {
			return nil, nil, err
		}

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
	varLexem, err := p.accept(models.LexLETTER, "variable name", "ожидалось название переменной")
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
	numLexem, err := p.accept(models.LexNUM, "number", "после { ожидалось значение степени - натуральное число")
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
