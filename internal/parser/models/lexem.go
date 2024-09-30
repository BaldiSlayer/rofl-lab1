package models

type Lexem struct {
	LexemType
	Str string
}

func (l Lexem) String() string {
	return l.Str
}

func (l Lexem) Type() LexemType {
	return l.LexemType
}

func NewEofLexem() Lexem {
	return Lexem{
		LexemType: LexEOF,
		Str:       "EOF",
	}
}

type LexemType int

const (
	LexVAR LexemType = iota
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
	LexEOF
)
