package trsparser

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/lexer"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"testing"
)

const (
	peano            string = "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T\n-------------S(x) = x+1\nf(x,y)=    x+2*y"
	wrongVar         string = "variables = x,y,z\n f(x, y) = f(x, z)\n-------f(x,y)     = xy"
	varError         string = "variables = x, y,\n f(x,y) = f(x,y)\n--------f(x,y) = x+y"
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
