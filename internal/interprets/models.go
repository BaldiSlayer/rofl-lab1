package interprets

type Interpretation struct {
	name string
	args []string
	monomials []Monomial
	constants []int
}

type Monomial struct {
	variable string
	coefficient int
	power int
}

type ParseError struct {
	llmMessage string
	summary string
}

func (e *ParseError) Error() string {
	return e.summary
}

func (e *ParseError) LlmMessage() string {
	return e.llmMessage
}
