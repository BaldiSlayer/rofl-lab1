package main;

import (
	"strconv"
	"fmt"
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
)

type Lexem struct{
	index, lex_type int
}

func lexer(a string) (res []Lexem){
	res = make([]Lexem, 0,len(a))
	variables := "variables"
	
	i := 0;
	for ; i < len(a); i++{
		switch a[i]{
			case ' ':
				continue
			case '\t':
				continue
			case '=':
				res =append(res,Lexem{i, lex_EQ})
			case ',':
				res=append(res,Lexem{i, lex_COMMA})
			case '+':
				res=append(res,Lexem{i, lex_ADD})
			case '*':
				res=append(res,Lexem{i, lex_MUL})
			case '{':
				res=append(res, Lexem{i, lex_LCB})
			case '}':
				res=append(res,Lexem{i, lex_RCB})
			case '(':
				res=append(res,Lexem{i, lex_LB})
			case ')':
				res=append(res,Lexem{i, lex_RB})
			default:
				if a[i] == '\n' || a[i] == '\r'{
					res=append(res, Lexem{i, lex_EOL})
					if i < len(a)-1 && (a[i] == '\n' || a[i] == '\r'){
						i++
					}
				}else if a[i] >= 'a' && a[i] <= 'z' || a[i] >= 'A' && a[i] <= 'Z'{
					if a[i] == 'v' && i + 8< len(a){
						t := true
						j := 0
						for ; j < 9; j++{
							if variables[j] != a[i+j]{
								t = false
								break
							}
						}
						if t{
							res=append(res,Lexem{i, lex_VAR})
							i += 8
						}else{
							for ; j >0; j--{
								res=append(res,Lexem{i, lex_LETTER})
								i++
							}
						}
					}else{
						res=append(res,Lexem{i, lex_LETTER})
					}
				}else if a[i] >= '0' && a[i] <= '9'{
					res=append(res,Lexem{i, lex_NUM})
					for; i< len(a) || (a[i] >= '0' && a[i] <= '9'); i++{} 
				}else{
					panic("unknown symbol at pos " + strconv.Itoa(i) + ":"+string(a[i]))
				}
		}
	}
	return res
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

func lexCheckOrPanic(l Lexem, Ltype int) bool{
	if l.lex_type != Ltype{
		panic("on index"+ strconv.Itoa(l.index)+ " expected "+ strconv.Itoa(Ltype)+ ", found " + strconv.Itoa(l.lex_type))
		return false
	}
	return true
}
func lexCheck(l Lexem, Ltype int) bool{
	return l.lex_type == Ltype
}

// <vars> ::= "variables" "=" <letters> <eol>
func TRS_parseVars(input string, m []Lexem, index *int){
	lexCheckOrPanic(m[*index], lex_VAR)
	*index++
	lexCheckOrPanic(m[*index], lex_EQ)
	*index++
	TRS_parseLetters(m, index)
	lexCheckOrPanic(m[*index], lex_EOL)
	*index++

}
//<letters> ::= <letter> <letters-tail>
func TRS_parseLetters(input string, m []Lexem, index *int){
	lexCheckOrPanic(m[*index], lex_LETTER)
	*index++
	TRS_parseLettersTail(m, index)
}

//<letters-tail> ::= "," <letter> <letters-tail> | ε
func TRS_parseLettersTail(m []Lexem, index *int){
	if lexCheck(m[*index], lex_COMMA){
		*index++
		lexCheckOrPanic(m[*index], lex_LETTER)
		*index++
		TRS_parseLettersTail(m, index)
	}
	
}


// <rules> ::= <rule> <eol> <rules-tail>
func TRS_parseRules(m []Lexem, index *int){}
func TRS_parseRulesTail(m []Lexem, index *int){}
func TRS_parseRule(m []Lexem, index *int){}
func TRS_parseTerm(m []Lexem, index *int){}
func TRS_parseArgs(m []Lexem, index *int){}
func TRS_parseTermsTail(m []Lexem, index *int){}

func parseTRS(m []Lexem) bool{
	return false
}

/*********************************************************************************/

func main(){
	test := "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T"
	a := lexer(test)
	for _, e := range a{
		fmt.Printf("%d %d\n", e.index, e.lex_type)
	}
}