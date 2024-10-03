package lexer

import (
	//"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"testing"
)

const (
	peano string = "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T\n-------------S(x) = x+1\nf(x,y)=    x+2*y"
)

func TestLexerWithPeano(t *testing.T) {
	//Peano grammar
	p := Lexer{text: peano}
	err := p.process()

	if err != nil {
		t.Errorf("Should not produce an error")
	}
}
