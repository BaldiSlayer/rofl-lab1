package trsparser

import (
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

<s> ::= <vars> <Rules>
<eol> ::= \n | \r | \r\n
<vars> ::= "variables" "=" <letters> <eol>
<letters> ::= <letter> <letters-tail>
<letters-tail> ::= "," <letter> <letters-tail> | ε
<Rules> ::= <rule> <eol> <Rules-tail>
<Rules-tail> ::= <rule> <eol> <Rules-tail> | ε
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

// Subexpression defines Model for Subexpression.
type Subexpression struct {
	Args *[]Subexpression

	// Letter represents variable or constructor
	Letter models.Lexem
}

type Parser struct {
	lexem []models.Lexem
	index int //index of syntax analyzing

	errors []*models.ParseError

	Model TRS
}

/*********************************************************************************/

/*
<s> ::= <vars> <Rules>

<vars> ::= "variables" "=" <letters> <eol>
<letters> ::= <letter> <letters-tail>
<letters-tail> ::= "," <letter> <letters-tail> | ε

<Rules> ::= <rule> <eol> <Rules-tail>
<Rules-tail> ::= <rule> <eol> <Rules-tail> | ε
<rule> ::= <term> "=" <term>
<term> ::= var | constructor <args>
<args> ::= ε | "(" <term> <terms-tail> ")"
<terms-tail> ::= "," <term> <terms-tail> | ε
*/

func (p *Parser) addRule() *Rule {
	i := len(p.Model.Rules)
	p.Model.Rules = append(p.Model.Rules, Rule{})
	return &p.Model.Rules[i]
}

func (p *Parser) addVariable(l models.Lexem) {
	p.Model.Variables = append(p.Model.Variables, l)
}

func (p *Parser) isVariable(l models.Lexem) bool {
	for _, e := range p.Model.Variables {
		if e.Str == l.Str {
			return true
		}
	}
	return false
}

func (p *Parser) lexCheck(l models.Lexem, Ltype models.LexemType) error {
	if l.LexemType != Ltype {
		switch l.LexemType {
		/*case models.LexLB:
			fallthrough
		case models.LexRB:
		*/default:
			return_str := models.GetLexemInterpretation(l.LexemType)
			if l.LexemType == models.LexNUM || l.LexemType == models.LexLETTER {
				return_str = l.Str
			}
			return &models.ParseError{
				LlmMessage: fmt.Sprintf("в строке %d TRS  на позиции %d ожидалось \"%s\", найдено \"%s\"",
					l.Line, l.Index, models.GetLexemInterpretation(Ltype), return_str),
				Message: fmt.Sprintf("at %d:%d expected %s, found %s",
					l.Line, l.Index, models.GetLexemInterpretation(Ltype), return_str),
			}
		}
	}
	return nil
}

// <vars> ::= "variables" "=" <letters> <eol>
func (p *Parser) parseVars() error {
	err := p.lexCheck(p.lexem[p.index], models.LexVAR)
	if err != nil {
		return models.Wrap(
			"variables definiton expected",
			"в начале TRS ожидалось перечисление переменных формата \"variables = x,y,z\"",
			err,
		)
	}
	p.index++
	err = p.lexCheck(p.lexem[p.index], models.LexEQ)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseLetters()
	if err != nil {
		return err
	}
	err = p.lexCheck(p.lexem[p.index], models.LexEOL)
	if err != nil {
		return err
	}
	p.index++
	return nil
}

// <letters> ::= <letter> <letters-tail>
func (p *Parser) parseLetters() error {
	err := p.lexCheck(p.lexem[p.index], models.LexLETTER)
	if err != nil {
		return models.Wrap(
			"at least one variable definition expected",
			"должна быть перечислена хоть одна переменная в списке переменных в первой строке",
			err,
		)
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
		err := p.lexCheck(p.lexem[p.index], models.LexLETTER)
		if err != nil {
			return err
		}
		p.addVariable(p.lexem[p.index])
		p.index++
		//p.parseLettersTail()
	}
	return nil
}

// <Rules> ::= <rule> <eol> <Rules-tail>
func (p *Parser) parseRules() error {
	err := p.parseRule()
	if err != nil {
		return err
	}
	err = p.lexCheck(p.lexem[p.index], models.LexEOL)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseRulesTail()
	return err
}

// <Rules-tail> ::= <rule> <eol> <Rules-tail> | ε
func (p *Parser) parseRulesTail() error {
	for p.lexem[p.index].LexemType == models.LexLETTER {
		err := p.parseRule()
		if err != nil {
			return err
		}
		err = p.lexCheck(p.lexem[p.index], models.LexEOL)
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

	err = p.lexCheck(p.lexem[p.index], models.LexEQ)
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
	err := p.lexCheck(p.lexem[p.index], models.LexLETTER)
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
		countInConstructor, ok := p.Model.Constructors[letter.Str]
		countVars := len(*subexpr_arr)
		if ok {
			if countInConstructor != countVars {
				return nil, &models.ParseError{
					LlmMessage: fmt.Sprintf("в строке %d несовпадение в количестве элементов конструктора %s: ожидалось %d аргументов, найдено %d аргументов",
						letter.Line, letter.Str, countInConstructor, countVars),
					Message: fmt.Sprintf("in line %d constructor mismatch %s: expect %d args, found %d args",
						letter.Line, letter.Str, countInConstructor, countVars),
				}
			}
		} else {
			p.Model.Constructors[letter.Str] = countVars
		}
		return &Subexpression{Args: subexpr_arr, Letter: letter}, nil
	} else { // variable
		if p.index < len(p.lexem) && p.lexem[p.index].LexemType == models.LexLB{
			return nil, &models.ParseError{
					LlmMessage: fmt.Sprintf("в строке %d позиции %d переменная %s использована как конструктор",
						letter.Line,letter.Index, letter.Str),
					Message: fmt.Sprintf("in line %d at %d var %s used as constructor",
						letter.Line,letter.Index, letter.Str),
				}
		}
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
		err = p.lexCheck(p.lexem[p.index], models.LexRB)
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

// <s> ::= <vars> <Rules>
func (p *Parser) parseTRS() error {
	p.index = 0
	p.Model = TRS{}
	p.Model.Constructors = make(map[string]int)
	p.errors = make([]*models.ParseError, 0)

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

func isSetIn(a, b *map[string]bool) (bool, string) {
	for element := range *b {
		if !((*a)[element]) {
			return false, element
		}
	}

	return true, ""
}

func (p *Parser) checkRules() error {
	for i, rule := range p.Model.Rules { // проверка корректности переменных
		left_var := make(map[string]bool)
		getVariablesFromExpr(&left_var, rule.Lhs)
		right_var := make(map[string]bool)
		getVariablesFromExpr(&right_var, rule.Rhs)
		isIn, contr := isSetIn(&left_var, &right_var)
		if !isIn {
			return &models.ParseError{
				LlmMessage: fmt.Sprintf("в правиле %d неправильно использованы переменные: %s не может быть использована в правой части",
					i+1, contr),
				Message: fmt.Sprintf("in rule %d var mismatch: %s cant be used",
					i+1, contr),
			}
		}
	}
	return nil
}

func ParseRules(arr []models.Lexem) (*TRS, []models.Lexem, error) {
	if len(arr) == 0 {
		return nil, arr, &models.ParseError{
			LlmMessage: "должно быть хотя бы одно правило переписывания",
			Message:    "need term rewrite rule",
		}
	}

	p := Parser{lexem: arr}

	err := p.parseTRS()
	if err != nil {
		return nil, arr, err
	}
	return &p.Model, p.lexem[p.index:], nil
}
