package models

type Lexem struct {
	LexemType
	Index int
	Str   string
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

func GetLexemInterpretation(l LexemType) string {
	names := [...]string{"variables", "=", "буква", ",", "*", "+", "{", "}", "(", ")", "число", "конец строки", "конец файла", "разделитель"}
	return names[l]
}

// <lexem> ::= "variables" | "=" | letter | "," | "*" | "{" | "}" | "(" | ")" | "+" | number | '\r' | \n | \r\n
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
	LexSEPARATOR
	//('-')* - separate TRS input and interpet input: can be deleted in the future
)
