package main

import (
	"testing"
	"./parseTRS"
)

func TestLexer(t *testing.T){
	text := "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T-------------S(x) = x+1\nf(x,y)=    x+2*y"
	//Peano grammar
	a, err := Lexer(text)
	
	if err != nil {
    	t.Errorf("Should not produce an error")
	}
}