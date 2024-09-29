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

//<lexem> ::= "variables" | "=" | letter | "," | "*" | "{" | "}" | "(" | ")" | "+" | number | '\r' | \n | \r\n
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

type Parser struct{
	text string
}

func (p Parser)Lexer() ([]Lexem, error) {
	runes := []rune(p.text)
	res := make([]Lexem, 0, len(runes))
	
	lexVARiables := []rune("variables")

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case ' ': // пробел и таб пропустить
			continue
		case '\t':
			continue
		case '-':
			res = append(res, Lexem{i, lexSEPARATOR})
			for i < len(runes) && runes[i] == '-'{
				i++
			}
		case '=':
			res = append(res, Lexem{i, lexEQ})
		case ',':
			res = append(res, Lexem{i, lexCOMMA})
		case '+':
			res = append(res, Lexem{i, lexADD})
		case '*':
			res = append(res, Lexem{i, lexMUL})
		case '{':
			res = append(res, Lexem{i, lexLCB})
		case '}':
			res = append(res, Lexem{i, lexRCB})
		case '(':
			res = append(res, Lexem{i, lexLB})
		case ')':
			res = append(res, Lexem{i, lexRB})
		default:
			if runes[i] == '\n' || runes[i] == '\r' { // если перевод строки(причем могут быть 2), добавить лексему перевод строки
				res = append(res, Lexem{i, lexEOL})
				if i < len(runes)-1 && (runes[i] == '\n' && runes[i+1] == '\r' || runes[i] == '\r' && runes[i+1] == '\n') {
					i++
				}
			} else if runes[i] >= 'a' && runes[i] <= 'z' || runes[i] >= 'A' && runes[i] <= 'Z' { // если встретилась буква
				if runes[i] == 'v' && i+8 < len(runes) { // проверяем на "variables"
					t := true
					j := 0
					for ; j < 9; j++ {
						if lexVARiables[j] != runes[i+j] {
							t = false
							break
						}
					}
					if t { // если найдено слово, добавляем и пропускаем
						res = append(res, Lexem{i, lexVAR})
						i += 8
					} else { // иначе добавляем букву 'v' и идем дальше посимвольно
						res = append(res, Lexem{i, lexLETTER})
					}
				} else { // если найденная буква не v, то добавляем букву
					res = append(res, Lexem{i, lexLETTER})
				}
			} else if runes[i] >= '0' && runes[i] <= '9' {
				res = append(res, Lexem{i, lexNUM})
				for i < len(runes) && (runes[i] >= '0' && runes[i] <= '9'){
					i++
				}
			} else {
				return nil, fmt.Errorf("unknown symbol at pos %d:%c", i, runes[i])
			}
		}
	}
	return res, nil
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

func lexCheck(l Lexem, Ltype int) (bool, error) {
	if l.lexType != Ltype {
		return false, fmt.Errorf("on index %d expected %d, found %d", l.index, Ltype, l.lexType)
	}
	return true, nil
}

// <vars> ::= "variables" "=" <letters> <eol>
func (p Parser) parseVars(/*input string,*/ m []Lexem, index *int) {
	lexCheck(m[*index], lexVAR)
	*index++
	lexCheck(m[*index], lexEQ)
	*index++
	p.parseLetters(m, index)
	lexCheck(m[*index], lexEOL)
	*index++

}

//<letters> ::= <letter> <letters-tail>
func (p Parser) parseLetters(/*input string,*/ m []Lexem, index *int) {
	lexCheck(m[*index], lexLETTER)
	*index++
	p.parseLettersTail(m, index)
}

//<letters-tail> ::= "," <letter> <letters-tail> | ε
func (p Parser) parseLettersTail(m []Lexem, index *int) {
	// вместо if оптимизировано с ипользованием цикла
	// для уменьшения глубины стека выполнения
	for m[*index].lexType == lexCOMMA {
		*index++
		lexCheck(m[*index], lexLETTER)
		*index++
		//p.parseLettersTail(m, index)
	}
}

// <rules> ::= <rule> <eol> <rules-tail>
func (p Parser) parseRules(m []Lexem, index *int)     {
	p.parseRule(m, index)
	lexCheck(m[*index], lexEOL)
	*index++
	p.parseRulesTail(m, index)
}
// <rules-tail> ::= <rule> <eol> <rules-tail> | ε
func (p Parser) parseRulesTail(m []Lexem, index *int) {
	for m[*index].lexType == lexLETTER{
		p.parseRule(m, index)
		lexCheck(m[*index], lexEOL)
		*index++
		//p.parseRulesTail(m, index)
	}
}
// <rule> ::= <term> "=" <term>
func (p Parser) parseRule(m []Lexem, index *int)      {
	p.parseTerm(m, index)
	lexCheck(m[*index], lexEQ)
	*index++
	p.parseTerm(m, index)
}
// <term> ::= var | constructor <args>
func (p Parser) parseTerm(m []Lexem, index *int)      {
	lexCheck(m[*index], lexLETTER)
	*index++
	if(true /*m[*index-1] not in vars*/){
		p.parseArgs(m, index)
	}
}
// <args> ::= ε | "(" <term> <terms-tail> ")"
func (p Parser) parseArgs(m []Lexem, index *int)      {
	if m[*index].lexType == lexLB{
		*index++
		p.parseTerm(m, index)
		p.parseTermsTail(m, index)
		lexCheck(m[*index], lexRB)
	}
}
// <terms-tail> ::= "," <term> <terms-tail> | ε
func (p Parser) parseTermsTail(m []Lexem, index *int) {
	for m[*index].lexType == lexCOMMA{
		*index++
		p.parseTerm(m, index)
		//p.parseTermsTail(m,index)
	}
}

// <s> ::= <vars> <rules>
func (p Parser) parseTRS(m []Lexem) bool {
	index := 0
	/*varList := */p.parseVars(m, &index)
	p.parseRules(m, &index/*, varList*/)
	return true
}
