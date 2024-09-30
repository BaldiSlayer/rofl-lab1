package trsinterprets

import "fmt"

type Interpretation struct {
	name      string
	args      []string
	monomials []Monomial
	constants []int
}

type Monomial struct {
	variable    string
	coefficient int
	power       int
}

type ParseError struct {
	llmMessage string
	message    string
}

func (child *ParseError) wrap(parent *ParseError) *ParseError {
	return &ParseError{
		llmMessage: fmt.Sprintf("%s: %s", parent.llmMessage, child.llmMessage),
		message:    fmt.Sprintf("%s: %s", parent.message, child.message),
	}
}

func (e *ParseError) Error() string {
	return e.message
}

func (e *ParseError) LlmMessage() string {
	return e.llmMessage
}
