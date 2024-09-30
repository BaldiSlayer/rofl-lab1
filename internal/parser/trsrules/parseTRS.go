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

type Parser struct {
	text  string
	lexem []models.Lexem
	index int //index of syntax analyzing
}

func (p *Parser) appendLex(index int, lexType models.LexemType, str string) {
	p.lexem = append(p.lexem, models.Lexem{/*index: index,*/ LexemType: lexType, Str: str})
}

func (p *Parser) Lexer() error {
	runes := []rune(p.text)
	p.lexem = make([]models.Lexem, 0, len(runes))

	lexVariables := []rune("variables")

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case ' ': // пробел и таб пропустить
			continue
		case '\t':
			continue
		case '-':
			p.appendLex(i, models.LexSEPARATOR, "-")
			for i < len(runes) && runes[i] == '-' {
				i++
			}
		case '=':
			p.appendLex(i, models.LexEQ, "=")
		case ',':
			p.appendLex(i, models.LexCOMMA, ",")
		case '+':
			p.appendLex(i, models.LexADD, "+")
		case '*':
			p.appendLex(i, models.LexMUL, "*")
		case '{':
			p.appendLex(i, models.LexLCB, "{")
		case '}':
			p.appendLex(i, models.LexRCB, "}")
		case '(':
			p.appendLex(i, models.LexLB, "(")
		case ')':
			p.appendLex(i, models.LexRB, ")")
		default:
			if runes[i] == '\n' || runes[i] == '\r' { // если перевод строки(причем могут быть 2), добавить лексему перевод строки
				p.appendLex(i, models.LexEOL, "\n")
				if i < len(runes)-1 && (runes[i] == '\n' && runes[i+1] == '\r' || runes[i] == '\r' && runes[i+1] == '\n') {
					i++
				}
			} else if runes[i] >= 'a' && runes[i] <= 'z' || runes[i] >= 'A' && runes[i] <= 'Z' { // если встретилась буква
				if runes[i] == 'v' && i+len(lexVariables) < len(runes) { // проверяем на "variables"
					t := true
					j := 0
					for ; j < 9; j++ {
						if lexVariables[j] != runes[i+j] {
							t = false
							break
						}
					}
					if t { // если найдено слово, добавляем и пропускаем
						p.appendLex(i, models.LexVAR, "variables")
						i += 8
					} else { // иначе добавляем букву 'v' и идем дальше посимвольно
						p.appendLex(i, models.LexLETTER, string(runes[i]))
					}
				} else { // если найденная буква не v, то добавляем букву
					p.appendLex(i, models.LexLETTER, string(runes[i]))
				}
			} else if runes[i] >= '0' && runes[i] <= '9' {
				p.appendLex(i, models.LexNUM, string(runes[i]))
				for i < len(runes) && (runes[i] >= '0' && runes[i] <= '9') {
					i++
				}
			} else {
				return fmt.Errorf("unknown symbol at pos %d:%c", i, runes[i])
			}
		}
	}
	return nil
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

func (p *Parser) isVariable(l models.Lexem) bool {
	return true
}

func lexCheck(l models.Lexem, Ltype models.LexemType) error {
	if l.LexemType != Ltype {
		return fmt.Errorf("on index %d expected %d, found %s", 0/*l.index*/, Ltype, l.Str)// todo: сделать подстановку str Ltype
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
	if err != nil{
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
	p.index++
	err = p.parseLettersTail()
	if err != nil{
		return err
	}
	return nil
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
		p.index++
		//p.parseLettersTail()
	}
	return nil
}

// <rules> ::= <rule> <eol> <rules-tail>
func (p *Parser) parseRules() error {
	err := p.parseRule()
	if err != nil{
		return err
	}
	err = lexCheck(p.lexem[p.index], models.LexEOL)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseRulesTail()
	if err != nil{
		return err
	}
	return nil
}

// <rules-tail> ::= <rule> <eol> <rules-tail> | ε
func (p *Parser) parseRulesTail() error {
	for p.lexem[p.index].LexemType == models.LexLETTER {
		err := p.parseRule()
		if err != nil{
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
	err := p.parseTerm()
	if err != nil{
		return err
	}
	err = lexCheck(p.lexem[p.index], models.LexEQ)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseTerm()
	if err != nil{
		return err
	}
	return nil
}

// <term> ::= var | constructor <args>
func (p *Parser) parseTerm() error {
	err := lexCheck(p.lexem[p.index], models.LexLETTER)
	if err != nil {
		return err
	}
	p.index++
	if !p.isVariable(p.lexem[p.index-1]) {
		err = p.parseArgs()
		if err != nil{
			return err
		}
	}
	return nil
}

// <args> ::= ε | "(" <term> <terms-tail> ")"
func (p *Parser) parseArgs() error {
	if p.lexem[p.index].LexemType == models.LexLB {
		p.index++
		err := p.parseTerm()
		if err != nil{
			return err
		}
		err = p.parseTermsTail()
		if err != nil{
			return err
		}
		err = lexCheck(p.lexem[p.index], models.LexRB)
		if err != nil {
			return err
		}
		p.index++
	}
	return nil
}

// <terms-tail> ::= "," <term> <terms-tail> | ε
func (p *Parser) parseTermsTail() error {
	for p.lexem[p.index].LexemType == models.LexCOMMA {
		p.index++
		err := p.parseTerm()
		if err != nil{
			return err
		}
		//p.parseTermsTail(m,index)
	}
	return nil
}

// <s> ::= <vars> <rules>
func (p *Parser) parseTRS() error {
	p.index = 0
	/*varList := */ err := p.parseVars()
	if err != nil {
		return err
	}
	err = p.parseRules( /*varList*/ )
	if err != nil {
		return err
	}
	return nil
}
