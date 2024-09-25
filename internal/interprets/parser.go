package interprets

import (
	"fmt"
	"strconv"

	"github.com/BaldiSlayer/rofl-lab1/internal/trs"
)

type Parser struct {
	stream stream
	ConstructorArity map[string]int
}

func NewParser(input chan trs.Lexem, constructorArity map[string]int) *Parser {
	return &Parser{
		stream:           stream{
			channel: input,
			ok:      false,
			val:     trs.Lexem{},
		},
		ConstructorArity: constructorArity,
	}
}

func (p *Parser) Parse() ([]Interpretation, error) {
	i, err := p.interprets()
	if err != nil {
		return nil, err
	}
	return i, nil
}

/*
<interprets> ::= <func-rule> <eol> <interprets> | <const-rule> <eol> <interprets> | ε
<const-rule> ::= <constructor> "=" number
<func-rule> ::= <constructor> "(" <letters> ")" "=" <monomial> <monomials-tail>
<monomials-tail> ::= "+" <monomial> <monomials-tail> | ε
<monomial> ::= number | <power-product> <power-products-tail>
<power-products-tail> ::= <power-product> <power-products-tail> | ε
<power-product> ::= <coeff> <var> <power>
<coeff> ::= ε | number "*"
<power> ::= ε | "{" number "}"
<letters> ::= letter <letters-tail>
<letters-tail> ::= "," letter <letters-tail> | ε
<eol> ::= \n | \r | \r\n
 */

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

func (p *Parser) interprets() ([]Interpretation, *ParseError) {
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

		constructor, err := p.accept(
			trs.LexLETTER,
			"constructor name",
			"ожидалось название конструктора",
		)
		if err != nil {
			return nil, err.wrap(&ParseError{
				llmMessage: "неверно задана интерпретация",
				message:    "wrong interpretation definition",
			})
		}

		interpret, err := p.constOrFuncRule(constructor.String())
		if err != nil {
			return nil, err.wrap(&ParseError{
				llmMessage: fmt.Sprintf("неверно задана интерпретация конструктора %s", constructor.String()),
				message:    "wrong interpretation definition",
			})
		}

		p.accept(trs.LexEOL, "EOL", "ожидался перенос строки")

		res = append(res, interpret)
	}
}

func (p *Parser) constOrFuncRule(name string) (Interpretation, *ParseError) {
	switch p.peek() {
	case trs.LexEQ:
		value, err := p.constRule()
		return Interpretation{
			name:      name,
			args:      []string{},
			monomials: []Monomial{},
			constants: []int{value},
		}, err
	case trs.LexLB:
		return p.funcRule(name)
	}

	got := p.stream.next()
	return Interpretation{}, &ParseError{
		llmMessage: fmt.Sprintf("ожидалось = или ( после названия конструктора, получено %s", got.String()),
		message:    fmt.Sprintf("expected = or (, got %s", got.String()),
	}
}

func (p *Parser) constRule() (int, *ParseError) {
	p.stream.next()
	lexem, err := p.accept(trs.LexNUM, "number", "ожидалось натуральное число после знака = в интерпретации константы")
	num, _ := strconv.Atoi(lexem.String())
	return num, err
}

// <func-rule> ::= <constructor> "(" <letters> ")" "=" <monomial> <monomials-tail>
func (p *Parser) funcRule(name string) (Interpretation, *ParseError) {
	p.stream.next()

	// TODO: check if name occures in args
	args, err := p.letters()
	if err != nil {
		return Interpretation{}, err
	}

	p.accept(trs.LexRB, ")", "TODO")

	return Interpretation{
		name:      name,
		args:      args,
		monomials: []Monomial{},
		constants: []int{},
	}, nil
}

func (p *Parser) letters() ([]string, *ParseError) {
	return nil, nil
}
