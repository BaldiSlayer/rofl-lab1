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

func (p *Parser) lexCheck(l models.Lexem, Ltype models.LexemType) *models.ParseError {
	if l.LexemType != Ltype {
		switch l.LexemType {
		/*case models.LexLB:
			fallthrough
		case models.LexRB:
			return fmt.Errorf("неправильная скобочная структура в строке %d TRS", p.lineIndex)
		*/default:
			return_str := models.GetLexemInterpretation(l.LexemType)
			if l.LexemType == models.LexNUM || l.LexemType == models.LexLETTER {
				return_str = l.Str
			}
			//return fmt.Errorf("в строке %d TRS  на позиции %d ожидалось \"%s\", найдено \"%s\"",
			//	l.Line+1, l.Index+1, models.GetLexemInterpretation(Ltype), return_str)
			return &models.ParseError{
				LlmMessage: fmt.Sprintf("в строке %d TRS  на позиции %d ожидалось \"%s\", найдено \"%s\"",
					l.Line+1, l.Index+1, models.GetLexemInterpretation(Ltype), return_str),
				Message: fmt.Sprintf("at %d:%d expected %s, found %s",
					l.Line+1, l.Index+1, models.GetLexemInterpretation(Ltype), return_str),
			}
		}
	}
	return nil
}

// <vars> ::= "variables" "=" <letters> <eol>
func (p *Parser) parseVars() *models.ParseError {
	err := p.lexCheck(p.lexem[p.index], models.LexVAR)
	if err != nil {
		//return fmt.Errorf()
		return &models.ParseError{
			LlmMessage: "в начале TRS ожидалось перечисление переменных формата \"variables = x,y,z\"",
			Message:    err.Error(),
		}
	}
	p.index++
	err = p.lexCheck(p.lexem[p.index], models.LexEQ)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseLetters()
	if err != nil {
		return err //errors.Join(fmt.Errorf("parseVars:"), err)
	}
	err = p.lexCheck(p.lexem[p.index], models.LexEOL)
	if err != nil {
		return err
	}
	p.index++
	return nil
}

// <letters> ::= <letter> <letters-tail>
func (p *Parser) parseLetters() *models.ParseError {
	err := p.lexCheck(p.lexem[p.index], models.LexLETTER)
	if err != nil {
		//return fmt.Errorf("должна быть перечислена хоть одна переменная в списке переменных в первой строке")
		return &models.ParseError{
			LlmMessage: "должна быть перечислена хоть одна переменная в списке переменных в первой строке",
			Message:    err.Error(),
		}
	}
	p.addVariable(p.lexem[p.index])
	p.index++
	err = p.parseLettersTail()
	return err
}

// <letters-tail> ::= "," <letter> <letters-tail> | ε
func (p *Parser) parseLettersTail() *models.ParseError {
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
func (p *Parser) parseRules() *models.ParseError {
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
func (p *Parser) parseRulesTail() *models.ParseError {
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
func (p *Parser) parseRule() *models.ParseError {
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
func (p *Parser) parseTerm() (*Subexpression, *models.ParseError) {
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
		return &Subexpression{Args: subexpr_arr, Letter: letter}, nil
	} else { // variable
		return &Subexpression{Args: nil, Letter: letter}, nil
	}
}

// <args> ::= ε | "(" <term> <terms-tail> ")"
func (p *Parser) parseArgs() (*[]Subexpression, *models.ParseError) {
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
func (p *Parser) parseTermsTail(arr *[]Subexpression) *models.ParseError {
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
func (p *Parser) parseTRS() *models.ParseError {
	p.index = 0
	p.Model = TRS{}
	err := p.parseVars()
	if err != nil {
		return err //, fmt.Errorf("Возможные способы решения: \n + Переменные должны состоять из одной буквы и разделены запятой"))
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

func (p *Parser) getConstructorsFromExpr(a Subexpression) *models.ParseError {
	if a.Args != nil {
		count, ok := p.Model.Constructors[a.Letter.Str]
		if !ok {
			p.Model.Constructors[a.Letter.Str] = len(*a.Args)
		} else {
			if count != len(*a.Args) {
				//return fmt.Errorf("несовпадение в количестве элементов конструктора %s: ожидалось %d переменных, найдено %d переменных", a.Letter.Str, count, len(*a.Args))
				return &models.ParseError{
					LlmMessage: fmt.Sprintf("несовпадение в количестве элементов конструктора %s: ожидалось %d переменных, найдено %d переменных",
						a.Letter.Str, count, len(*a.Args)),
					Message: fmt.Sprintf("constructor mismatch %s: expect %d vars, found %d vars",
						a.Letter.Str, count, len(*a.Args)),
				}
			}
		}
		for _, e := range *a.Args {
			p.getConstructorsFromExpr(e)
		}
	}
	return nil
}

func isSetIn(a, b *map[string]bool) (bool, string) {
	for element := range *b {
		if !((*a)[element]) {
			return false, element
		}
	}

	return true, ""
}

func (p *Parser) checkRules() *models.ParseError {
	for i, rule := range p.Model.Rules { // проверка корректности переменных
		left_var := make(map[string]bool)
		getVariablesFromExpr(&left_var, rule.Lhs)
		right_var := make(map[string]bool)
		getVariablesFromExpr(&right_var, rule.Rhs)
		isIn, contr := isSetIn(&left_var, &right_var)
		if !isIn {
			//return fmt.Errorf("в правиле %d неправильно использованы переменные: %s не может быть использована в правой части", i+1, contr)
			return &models.ParseError{
				LlmMessage: fmt.Sprintf("в правиле %d неправильно использованы переменные: %s не может быть использована в правой части",
					i+1, contr),
				Message: fmt.Sprintf("in rule %d var mismatch: %s cant be used",
					i+1, contr),
			}
		}
	}

	p.Model.Constructors = make(map[string]int)

	for i, rule := range p.Model.Rules {
		err := p.getConstructorsFromExpr(rule.Lhs)
		if err != nil {
			//return errors.Join(fmt.Errorf("в левой части правила %d ", i+1), err)
			return err.Wrap(&models.ParseError{
				LlmMessage: fmt.Sprintf("в левой части правила %d ", i+1),
				Message:    fmt.Sprintf("in left part of rule %d ", i+1),
			})
		}
		err = p.getConstructorsFromExpr(rule.Rhs)
		if err != nil {
			//return errors.Join(fmt.Errorf("в правой части правила %d ", i+1), err)
			return err.Wrap(&models.ParseError{
				LlmMessage: fmt.Sprintf("в правой части правила %d ", i+1),
				Message:    fmt.Sprintf("in right part of rule %d ", i+1),
			})
		}
	}

	return nil
}

func ParseRules(arr []models.Lexem) (*TRS, []models.Lexem, *models.ParseError) {
	if len(arr) == 0 {
		//return nil, arr, fmt.Errorf("должно быть хотя бы одно правило переписывания")
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
