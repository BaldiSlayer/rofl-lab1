package trsparser

import (
	"testing"
)

func TestLexerWithPeano(t *testing.T){
	text := "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T-------------S(x) = x+1\nf(x,y)=    x+2*y"
	//Peano grammar
	p := Parser{text, [], 0}
	err := p.Lexer()
	
	if err != nil {
    	t.Errorf("Should not produce an error")
	}
}