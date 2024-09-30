package trsparser

import (
	"fmt"
	//"errors"
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

// <lexem> ::= "variables" | "=" | letter | "," | "*" | "{" | "}" | "(" | ")" | "+" | number | '\r' | \n | \r\n
const (
	lexVAR int = iota
	lexEQ
	lexLETTER
	lexCOMMA
	lexMUL
	lexADD
	lexLCB
	lexRCB
	lexLB
	lexRB
	lexNUM
	lexEOL
	lexSEPARATOR //('-')* - separate TRS input and interpet input: can be deleted in the future
)

type Lexem struct {
	index, lexType int
}

type Parser struct {
	text  string
	lexem []Lexem
	index int //index of syntax analyzing
}

func (p *Parser) appendLex(index, lexType int) {
	p.lexem = append(p.lexem, Lexem{index, lexType})
}

func (p *Parser) Lexer() error {
	runes := []rune(p.text)
	p.lexem = make([]Lexem, 0, len(runes))

	lexVariables := []rune("variables")

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case ' ': // пробел и таб пропустить
			continue
		case '\t':
			continue
		case '-':
			p.appendLex(i, lexSEPARATOR)
			for i < len(runes) && runes[i] == '-' {
				i++
			}
		case '=':
			p.appendLex(i, lexEQ)
		case ',':
			p.appendLex(i, lexCOMMA)
		case '+':
			p.appendLex(i, lexADD)
		case '*':
			p.appendLex(i, lexMUL)
		case '{':
			p.appendLex(i, lexLCB)
		case '}':
			p.appendLex(i, lexRCB)
		case '(':
			p.appendLex(i, lexLB)
		case ')':
			p.appendLex(i, lexRB)
		default:
			if runes[i] == '\n' || runes[i] == '\r' { // если перевод строки(причем могут быть 2), добавить лексему перевод строки
				p.appendLex(i, lexEOL)
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
						p.appendLex(i, lexVAR)
						i += 8
					} else { // иначе добавляем букву 'v' и идем дальше посимвольно
						p.appendLex(i, lexLETTER)
					}
				} else { // если найденная буква не v, то добавляем букву
					p.appendLex(i, lexLETTER)
				}
			} else if runes[i] >= '0' && runes[i] <= '9' {
				p.appendLex(i, lexNUM)
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

func (p *Parser) isVariable(l Lexem) bool {
	return true
}

func lexCheck(l Lexem, Ltype int) error {
	if l.lexType != Ltype {
		return fmt.Errorf("on index %d expected %d, found %d", l.index, Ltype, l.lexType)
	}
	return nil
}

// <vars> ::= "variables" "=" <letters> <eol>
func (p *Parser) parseVars() error {
	err := lexCheck(p.lexem[p.index], lexVAR)
	if err != nil {
		return err
	}
	p.index++
	err = lexCheck(p.lexem[p.index], lexEQ)
	if err != nil {
		return err
	}
	p.index++
	p.parseLetters()
	err = lexCheck(p.lexem[p.index], lexEOL)
	if err != nil {
		return err
	}
	p.index++
	return nil
}

// <letters> ::= <letter> <letters-tail>
func (p *Parser) parseLetters() error {
	err := lexCheck(p.lexem[p.index], lexLETTER)
	if err != nil {
		return err
	}
	p.index++
	p.parseLettersTail()
	return nil
}

// <letters-tail> ::= "," <letter> <letters-tail> | ε
func (p *Parser) parseLettersTail() error {
	// вместо if оптимизировано с ипользованием цикла
	// для уменьшения глубины стека выполнения
	for p.lexem[p.index].lexType == lexCOMMA {
		p.index++
		err := lexCheck(p.lexem[p.index], lexLETTER)
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
	p.parseRule()
	err := lexCheck(p.lexem[p.index], lexEOL)
	if err != nil {
		return err
	}
	p.index++
	p.parseRulesTail()
	return nil
}

// <rules-tail> ::= <rule> <eol> <rules-tail> | ε
func (p *Parser) parseRulesTail() error {
	for p.lexem[p.index].lexType == lexLETTER {
		p.parseRule()
		err := lexCheck(p.lexem[p.index], lexEOL)
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
	p.parseTerm()
	err := lexCheck(p.lexem[p.index], lexEQ)
	if err != nil {
		return err
	}
	p.index++
	p.parseTerm()
	return nil
}

// <term> ::= var | constructor <args>
func (p *Parser) parseTerm() error {
	err := lexCheck(p.lexem[p.index], lexLETTER)
	if err != nil {
		return err
	}
	p.index++
	if !p.isVariable(p.lexem[p.index-1]) {
		p.parseArgs()
	}
	return nil
}

// <args> ::= ε | "(" <term> <terms-tail> ")"
func (p *Parser) parseArgs() error {
	if p.lexem[p.index].lexType == lexLB {
		p.index++
		p.parseTerm()
		p.parseTermsTail()
		err := lexCheck(p.lexem[p.index], lexRB)
		if err != nil {
			return err
		}
		p.index++
	}
	return nil
}

// <terms-tail> ::= "," <term> <terms-tail> | ε
func (p *Parser) parseTermsTail() error {
	for p.lexem[p.index].lexType == lexCOMMA {
		p.index++
		p.parseTerm()
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
