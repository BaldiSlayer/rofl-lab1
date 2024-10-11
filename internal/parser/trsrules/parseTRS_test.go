package trsparser

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/lexer"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"testing"
)

const (
	peano            string = "variables = x,y,z\n\n\n\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T\n-------------S(x) = x+1\nf(x,y)=    x+2*y"
	wrongConstructor        = "variables = x,y,z\n f(x,y) = f(x)\n----------f(x,y) = x"
)

func TestParserWithPeano(t *testing.T) {
	l := lexer.Lexer{Text: peano}
	err := l.Process()
	if err != nil {
		t.Error(err)
	}
	_, lex_tail, err1 := ParseRules(l.Lexem)
	if err1 != nil {
		t.Error(err1)
	}
	if lex_tail[0].LexemType != models.LexSEPARATOR {
		t.Errorf("Expected separator, but find %d lexem", lex_tail[0].LexemType)
	}
}
func TestParserWithWrongVar(t *testing.T) {
	wrongVar := "variables = x,y,z\n f(x, y) = f(x, z)\n-------f(x,y)     = xy"
	l := lexer.Lexer{Text: wrongVar}
	err := l.Process()
	if err != nil {
		t.Error(err)
	}

	_, _, err1 := ParseRules(l.Lexem)
	if err1.Error() != "in rule 1 var mismatch: z cant be used" {
		t.Error(err1)
	}
}
func TestParserWithVarError(t *testing.T) {
	varError := "variables = x, y,\n f(x,y) = f(x,y)\n--------f(x,y) = x+y"
	l := lexer.Lexer{Text: varError}
	err := l.Process()
	if err != nil {
		t.Error(err)
	}

	_, _, err1 := ParseRules(l.Lexem)
	if err1 != nil && err1.Error() != "at 1:18 expected буква, found конец строки" {
		t.Error(err1)
	}
}
func TestParserWithConstructorError(t *testing.T) {
	l := lexer.Lexer{Text: wrongConstructor}
	err := l.Process()
	if err != nil {
		t.Error(err)
	}

	_, _, err1 := ParseRules(l.Lexem)
	if err1.Error() != "in line 2 constructor mismatch f: expect 2 args, found 1 args" {
		t.Error(err1)
	}
}

func TestParserWithVarAsConstructor(t *testing.T) {
	l := lexer.Lexer{Text: "variables=x, y\nx(y)=y\n----x(y) = y+1\n"}
	err := l.Process()
	if err != nil {
		t.Error(err)
	}

	_, _, err1 := ParseRules(l.Lexem)
	if err1.Error() != "in line 2 at 1 var x used as constructor" {
		t.Error(err1)
	}
}

func TestMind(t *testing.T) {
	input := `variables=x
f(x) = 
-------`
	l := lexer.Lexer{Text: input}
	err := l.Process()
	if err != nil {
		t.Error(err)
	}

	_, _, err1 := ParseRules(l.Lexem)
	//if err1 != nil {
	//	t.Error(err1)
	//}
	if err1 == nil {
		t.Error("должен бросать ошибку о неправильной лексеме в второй строчке")
	}
}
