package trs

type Lexem interface {
	String() string
	Type() int
}

const (
	LexVAR int = iota
	LexEQ
	LexLETTER
	LexCOMMA
	LexMUL
	LexADD
	LexLCB
	LexRCB
	LexLB
	LexRB
	LexNUM
	LexEOL
)
