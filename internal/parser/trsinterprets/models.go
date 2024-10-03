package trsinterprets

import "fmt"

type Interpretation struct {
	name      string
	args      []string
	monomials []Monomial
}

// NOTE: one of {constant, factors}
type Monomial struct {
	constant *int
	factors *[]Factor
}

func NewConstantMonomial(v int) Monomial {
	return Monomial{
		constant: &v,
		factors:  nil,
	}
}

func NewProductMonomial(factors []Factor) Monomial {
	return Monomial{
		constant: nil,
		factors:  &factors,
	}
}

type Factor struct {
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
