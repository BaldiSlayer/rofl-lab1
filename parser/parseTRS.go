package main

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
	lex_VAR int = iota
	lex_EQ
	lex_LETTER
	lex_COMMA
	lex_MUL
	lex_ADD
	lex_LCB
	lex_RCB
	lex_LB
	lex_RB
	lex_NUM
	lex_EOL
	lex_SEPARATOR //('-')* - separate TRS input and interpet input: can be deleted in the future
)

type Lexem struct {
	index, lex_type int
}

func Lexer(text string) ([]Lexem, error) {
	runes := []rune(text)
	res := make([]Lexem, 0, len(runes))
	
	lex_variables := []rune("variables")

	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case ' ': // пробел и таб пропустить
			continue
		case '\t':
			continue
		case '-':
			res = append(res, Lexem{i, lex_SEPARATOR})
			for i < len(runes) && runes[i] == '-'{
				i++
			}
		case '=':
			res = append(res, Lexem{i, lex_EQ})
		case ',':
			res = append(res, Lexem{i, lex_COMMA})
		case '+':
			res = append(res, Lexem{i, lex_ADD})
		case '*':
			res = append(res, Lexem{i, lex_MUL})
		case '{':
			res = append(res, Lexem{i, lex_LCB})
		case '}':
			res = append(res, Lexem{i, lex_RCB})
		case '(':
			res = append(res, Lexem{i, lex_LB})
		case ')':
			res = append(res, Lexem{i, lex_RB})
		default:
			if runes[i] == '\n' || runes[i] == '\r' { // если перевод строки(причем могут быть 2), добавить лексему перевод строки
				res = append(res, Lexem{i, lex_EOL})
				if i < len(runes)-1 && (runes[i] == '\n' && runes[i+1] == '\r' || runes[i] == '\r' && runes[i+1] == '\n') {
					i++
				}
			} else if runes[i] >= 'a' && runes[i] <= 'z' || runes[i] >= 'A' && runes[i] <= 'Z' { // если встретилась буква
				if runes[i] == 'v' && i+8 < len(runes) { // проверяем на "variables"
					t := true
					j := 0
					for ; j < 9; j++ {
						if lex_variables[j] != runes[i+j] {
							t = false
							break
						}
					}
					if t { // если найдено слово, добавляем и пропускаем
						res = append(res, Lexem{i, lex_VAR})
						i += 8
					} else { // иначе добавляем букву 'v' и идем дальше посимвольно
						res = append(res, Lexem{i, lex_LETTER})
					}
				} else { // если найденная буква не v, то добавляем букву
					res = append(res, Lexem{i, lex_LETTER})
				}
			} else if runes[i] >= '0' && runes[i] <= '9' {
				res = append(res, Lexem{i, lex_NUM})
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
	if l.lex_type != Ltype {
		return false, fmt.Errorf("on index %d expected %d, found %d", l.index, Ltype, l.lex_type)
	}
	return true, nil
}

// <vars> ::= "variables" "=" <letters> <eol>
func TRS_parseVars(/*input string,*/ m []Lexem, index *int) {
	lexCheck(m[*index], lex_VAR)
	*index++
	lexCheck(m[*index], lex_EQ)
	*index++
	TRS_parseLetters(m, index)
	lexCheck(m[*index], lex_EOL)
	*index++

}

//<letters> ::= <letter> <letters-tail>
func TRS_parseLetters(/*input string,*/ m []Lexem, index *int) {
	lexCheck(m[*index], lex_LETTER)
	*index++
	TRS_parseLettersTail(m, index)
}

//<letters-tail> ::= "," <letter> <letters-tail> | ε
func TRS_parseLettersTail(m []Lexem, index *int) {
	if m[*index].lex_type == lex_COMMA {
		*index++
		lexCheck(m[*index], lex_LETTER)
		*index++
		TRS_parseLettersTail(m, index)
	}
}

// <rules> ::= <rule> <eol> <rules-tail>
func TRS_parseRules(m []Lexem, index *int)     {
	TRS_parseRule(m, index)
	lexCheck(m[*index], lex_EOL)
	*index++
	TRS_parseRulesTail(m, index)
}
// <rules-tail> ::= <rule> <eol> <rules-tail> | ε
func TRS_parseRulesTail(m []Lexem, index *int) {
	if m[*index].lex_type == lex_LETTER{
		TRS_parseRule(m, index);
		lexCheck(m[*index], lex_EOL)
		*index++
		TRS_parseRulesTail(m, index)
	}
}
func TRS_parseRule(m []Lexem, index *int)      {}
func TRS_parseTerm(m []Lexem, index *int)      {}
func TRS_parseArgs(m []Lexem, index *int)      {}
func TRS_parseTermsTail(m []Lexem, index *int) {}

func parseTRS(m []Lexem) bool {
	return false
}
