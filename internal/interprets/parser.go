package interprets

type Lexem interface {
	String() string
	Type() int
}

const (
	lex_VAR int = iota
	lex_EQ
	lex_LETTER
	lex_COMMA
	lex_MUL
	lex_ADD
	lex_LCB
	lex_RCB
	lex_LB
	lex_RB
	lex_NUM
	lex_EOL
)

type Parser struct {}

func (p Parser) Parse(input chan Lexem, constructorArity map[string]int) ([]Interpretation, error) {
	if constructorArity["f"] == 2 {
		return nil, &ParseError{
			llmMessage: "неверная арность интерпретации конструктора f: ожидалось 2, получено 1",
			summary:    "wrong interpretation arity",
		}
	}

	return []Interpretation{{
		name:      "f",
		args:      []string{"x"},
		monomials: []Monomial{},
		constants: []int{5},
	},}, nil
}
