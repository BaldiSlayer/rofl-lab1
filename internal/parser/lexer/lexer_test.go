package lexer

import (
	"testing"
)

const (
	peano string = "variables = x,y,z\n f(x,S(y)) = S(f(x,y)) \n\r f(x, T) = T\n-------------S(x) = x+1\nf(x,y)=    x+2*y"
)

func TestLexerWithPeano(t *testing.T) {
	//Peano grammar
	p := Lexer{Text: peano}
	err := p.Process()

	if err != nil {
		t.Errorf("Should not produce an error")
	}
}

func TestUnknownSymbol(t *testing.T) {
	p := Lexer{Text: "шалаш"}
	err := p.Process()
	if err == nil {
		t.Errorf("Should produce error %s", ">Неизвестный символ в строке 1, позиции 1: ш")
	}
}

func TestUnknownSymbolNotInFirstLine(t *testing.T) {
	p := Lexer{Text: "variables=x\nШ(x)=x\n------Ш(x)=x+1\n"}
	err := p.Process()
	if err == nil {
		t.Errorf("Should produce error %s", ">Неизвестный символ в строке 2, позиции 1: ш")
	}
}
