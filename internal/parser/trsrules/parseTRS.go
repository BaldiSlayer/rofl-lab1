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

type TRS struct {
	variables []models.Lexem
}

type Parser struct {
	text  string
	lexem []models.Lexem
	index int //index of syntax analyzing

	model TRS
}

func (p *Parser) appendLex(index int, lexType models.LexemType, str string) {
	p.lexem = append(p.lexem, models.Lexem{
		/*index: index,*/
		LexemType: lexType,
		Str:       str,
	})
}

func isLetter(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
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
			} else if isLetter(runes[i]) { // если встретилась буква
				if runes[i] == 'v' && i+len(lexVariables) < len(runes) { // проверяем на "variables"
					t := true
					j := 0
					for ; j < len(lexVariables); j++ {
						if lexVariables[j] != runes[i+j] {
							t = false
							break
						}
					}
					if t { // если найдено слово, добавляем и пропускаем
						p.appendLex(i, models.LexVAR, "variables")
						i += len(lexVariables) - 1
					} else { // иначе добавляем букву 'v' и идем дальше посимвольно
						p.appendLex(i, models.LexLETTER, string(runes[i]))
					}
				} else { // если найденная буква не v, то добавляем букву
					p.appendLex(i, models.LexLETTER, string(runes[i]))
				}
			} else if isDigit(runes[i]) {
				start_index := i
				for i+1 < len(runes) && isDigit(runes[i+1]) {
					i++
				}
				p.appendLex(i, models.LexNUM, string(runes[start_index:i]))
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

func (p *Parser) addVariable(l models.Lexem) {
	p.model.variables = append(p.model.variables, l)
}

func (p *Parser) isVariable(l models.Lexem) bool {
	for _, e := range p.model.variables {
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
	err := p.parseTerm()
	if err != nil {
		return err
	}
	err = lexCheck(p.lexem[p.index], models.LexEQ)
	if err != nil {
		return err
	}
	p.index++
	err = p.parseTerm()
	return err
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
		if err != nil {
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
		if err != nil {
			return err
		}
		err = p.parseTermsTail()
		if err != nil {
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
		if err != nil {
			return err
		}
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
	return err
}
