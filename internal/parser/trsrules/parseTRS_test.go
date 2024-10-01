package trsparser

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"testing"
)

const (
	peano string = "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T\n-------------S(x) = x+1\nf(x,y)=    x+2*y"
)

func TestLexerWithPeano(t *testing.T) {
	//Peano grammar
	p := Parser{text: peano}
	err := p.Lexer()

	if err != nil {
		t.Errorf("Should not produce an error")
	}
}

func TestParserWithPeano(t *testing.T) {
	p := Parser{text: peano}
	err := p.Lexer()
	if err != nil {
		t.Error(err)
	}
	err = p.parseTRS()
	if err != nil {
		t.Error(err)
	}
	if p.lexem[p.index].LexemType != models.LexSEPARATOR {
		t.Errorf("Expected separator, but find %d lexem", p.lexem[p.index].LexemType)
	}
}