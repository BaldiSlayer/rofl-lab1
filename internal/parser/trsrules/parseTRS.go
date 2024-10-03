package trsparser

import (
	"errors"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
)

/*
<lexem> ::= "variables" | "=" | letter | "," | "*" | "{" | "}" | "(" | ")" | "+" | number | '\r' | \n | \r\n

<variables> = "variables"
<eq> = '='
<letter> = буква
<comma> = ','
<mul> = '*'
<add> = '+'
<lcb> = '{'
<rcb> = '}'
<lb> = '('
<rb> = ')'
<num> = number
<eol> = '\n' | '\r' | "\n\r" | "\r\n"

grammatic

<s> ::= <vars> <rules>
<eol> ::= \n | \r | \r\n
<vars> ::= "variables" "=" <letters> <eol>
<letters> ::= <letter> <letters-tail>
<letters-tail> ::= "," <letter> <letters-tail> | ε
<rules> ::= <rule> <eol> <rules-tail>
<rules-tail> ::= <rule> <eol> <rules-tail> | ε
<rule> ::= <term> "=" <term>
<term> ::= var | constructor <args>
<args> ::= ε | "(" <term> <terms-tail> ")"
<terms-tail> ::= "," <term> <terms-tail> | ε
*/

type TRS struct {
	Variables []models.Lexem
	Rules     []Rule

	Constructors map[string]int
}

type Rule struct {
	Lhs Subexpression
	Rhs Subexpression
}

// Subexpression defines model for Subexpression.
type Subexpression struct {
	Args *[]Subexpression

	// Letter represents variable or constructor
	Letter models.Lexem
}

type Parser struct {
	lexem []models.Lexem
	index int //index of syntax analyzing

	model TRS
}

/*********************************************************************************/

/*
<s> ::= <vars> <rules>

<vars> ::= "variables" "=" <letters> <eol>
<letters> ::= <letter> <letters-tail>
<letters-tail> ::= "," <letter> <letters-tail> | ε

<rules> ::= <rule> <eol> <rules-tail>
<rules-tail> ::= <rule> <eol> <rules-tail> | ε
<rule> ::= <term> "=" <term>
<term> ::= var | constructor <args>
<args> ::= ε | "(" <term> <terms-tail> ")"
<terms-tail> ::= "," <term> <terms-tail> | ε
*/

func (p *Parser) addRule() *Rule {
	i := len(p.model.Rules)
	p.model.Rules = append(p.model.Rules, Rule{})
	return &p.model.Rules[i]
}

func (p *Parser) addVariable(l models.Lexem) {
	p.model.Variables = append(p.model.Variables, l)
}

func (p *Parser) isVariable(l models.Lexem) bool {
	for _, e := range p.model.Variables {
		if e.Str == l.Str {
			return true
		}
	}
	return false
}

func lexCheck(l models.Lexem, Ltype models.LexemType) error {
	if l.LexemType != Ltype {
		switch l.LexemType {
		case models.LexLB:
			fallthrough
		case models.LexRB:
			return fmt.Errorf("неправильная скобочная структура")
		default:
			return fmt.Errorf("on index %d expected %d, found %s", 0 /*l.index*/, Ltype, l.Str) // todo: сделать подстановку str Ltype
		}

	}
	return nil
}

// <vars> ::= "variables" "=" <letters> <eol>
func (p *Parser) parseVars() error {
	err := lexCheck(p.lexem[p.index], models.LexVAR)
	if err != nil {
		return err
	}
	p.index++
	err = lexCheck(p.lexem[p.index], models.LexEQ)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseLetters()
	if err != nil {
		return err
	}
	err = lexCheck(p.lexem[p.index], models.LexEOL)
	if err != nil {
		return err
	}
	p.index++
	return nil
}

// <letters> ::= <letter> <letters-tail>
func (p *Parser) parseLetters() error {
	err := lexCheck(p.lexem[p.index], models.LexLETTER)
	if err != nil {
		return err
	}
	p.addVariable(p.lexem[p.index])
	p.index++
	err = p.parseLettersTail()
	return err
}

// <letters-tail> ::= "," <letter> <letters-tail> | ε
func (p *Parser) parseLettersTail() error {
	// вместо if оптимизировано с ипользованием цикла
	// для уменьшения глубины стека выполнения
	for p.lexem[p.index].LexemType == models.LexCOMMA {
		p.index++
		err := lexCheck(p.lexem[p.index], models.LexLETTER)
		if err != nil {
			return err
		}
		p.addVariable(p.lexem[p.index])
		p.index++
		//p.parseLettersTail()
	}
	return nil
}

// <rules> ::= <rule> <eol> <rules-tail>
func (p *Parser) parseRules() error {
	err := p.parseRule()
	if err != nil {
		return err
	}
	err = lexCheck(p.lexem[p.index], models.LexEOL)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseRulesTail()
	return err
}

// <rules-tail> ::= <rule> <eol> <rules-tail> | ε
func (p *Parser) parseRulesTail() error {
	for p.lexem[p.index].LexemType == models.LexLETTER {
		err := p.parseRule()
		if err != nil {
			return err
		}
		err = lexCheck(p.lexem[p.index], models.LexEOL)
		if err != nil {
			return err
		}
		p.index++
		//p.parseRulesTail()
	}
	return nil
}

// <rule> ::= <term> "=" <term>
func (p *Parser) parseRule() error {
	r := p.addRule() // return *Rule

	subexp, err := p.parseTerm()
	if err != nil {
		return err
	}
	r.Lhs = *subexp

	err = lexCheck(p.lexem[p.index], models.LexEQ)
	if err != nil {
		return err
	}
	p.index++
	subexp, err = p.parseTerm()
	if err != nil {
		return err
	}
	r.Rhs = *subexp
	return nil
}

// <term> ::= var | constructor <args>
func (p *Parser) parseTerm() (*Subexpression, error) {
	err := lexCheck(p.lexem[p.index], models.LexLETTER)
	if err != nil {
		return nil, err
	}

	letter := p.lexem[p.index]

	p.index++
	if !p.isVariable(p.lexem[p.index-1]) { //constructor
		subexpr_arr, err1 := p.parseArgs()
		if err1 != nil {
			return nil, err1
		}
		return &Subexpression{Args: subexpr_arr, Letter: letter}, nil
	} else { // variable
		return &Subexpression{Args: nil, Letter: letter}, nil
	}
}

// <args> ::= ε | "(" <term> <terms-tail> ")"
func (p *Parser) parseArgs() (*[]Subexpression, error) {
	subexpr_arr := make([]Subexpression, 0)
	if p.lexem[p.index].LexemType == models.LexLB {
		p.index++
		subexpr, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		subexpr_arr = append(subexpr_arr, *subexpr)

		err = p.parseTermsTail(&subexpr_arr)
		if err != nil {
			return nil, err
		}
		err = lexCheck(p.lexem[p.index], models.LexRB)
		if err != nil {
			return nil, err
		}
		p.index++
	}
	return &subexpr_arr, nil
}

// <terms-tail> ::= "," <term> <terms-tail> | ε
func (p *Parser) parseTermsTail(arr *[]Subexpression) error {
	for p.lexem[p.index].LexemType == models.LexCOMMA {
		p.index++
		se, err := p.parseTerm()
		if err != nil {
			return err
		}
		*arr = append(*arr, *se)
		//p.parseTermsTail(m,index)
	}
	return nil
}

// <s> ::= <vars> <rules>
func (p *Parser) parseTRS() error {
	p.index = 0
	p.model = TRS{}
	err := p.parseVars()
	if err != nil {
		return err
	}
	err = p.parseRules()
	if err != nil {
		return err
	}

	return p.checkRules()
}

func getVariablesFromExpr(var_set *map[string]bool, a Subexpression) {
	if a.Args == nil {
		(*var_set)[a.Letter.Str] = true
	} else {
		for _, e := range *a.Args {
			getVariablesFromExpr(var_set, e)
		}
	}
}

func (p *Parser) getConstructorsFromExpr(a Subexpression) error {
	if a.Args != nil {
		count, ok := p.model.Constructors[a.Letter.Str]
		if !ok {
			p.model.Constructors[a.Letter.Str] = len(*a.Args)
		} else {
			if count != len(*a.Args) {
				return fmt.Errorf("несовпадение в количестве элементов конструктора %s: ожидалось %d, найдено %d", a.Letter.Str, count, len(*a.Args))
			}
		}
		for _, e := range *a.Args {
			p.getConstructorsFromExpr(e)
		}
	}
	return nil
}

func isSetIn(a, b *map[string]bool) bool {
	if len(*a) < len(*b) {
		return false
	}
	for element := range *b {
		if !((*a)[element]) {
			return false
		}
	}

	return true
}

func (p *Parser) checkRules() error {
	for i, rule := range p.model.Rules { // проверка корректности переменных
		left_var := make(map[string]bool)
		getVariablesFromExpr(&left_var, rule.Lhs)
		right_var := make(map[string]bool)
		getVariablesFromExpr(&right_var, rule.Rhs)
		if !isSetIn(&left_var, &right_var) {
			return fmt.Errorf("в правиле %d неправильно использованы переменные", i+1)
		}
	}

	p.model.Constructors = make(map[string]int)

	for i, rule := range p.model.Rules {
		err := p.getConstructorsFromExpr(rule.Lhs)
		if err != nil {
			return errors.Join(fmt.Errorf("в левой части правила %d ", i), err)
		}
		err = p.getConstructorsFromExpr(rule.Rhs)
		if err != nil {
			return errors.Join(fmt.Errorf("в правой части правила %d ", i), err)
		}
	}

	return nil
}

func ParseRules(arr []models.Lexem) (*TRS, []models.Lexem, error) {
	p := Parser{lexem: arr}

	err := p.parseTRS()
	if err != nil {
		return nil, arr, err
	}
	return &p.model, p.lexem[p.index:], nil
}
