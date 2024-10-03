package trsinterprets

import "fmt"

type Interpretation struct {
	Name      string
	Args      []string
	Monomials []Monomial
}

// NOTE: one of {constant, factors}
type Monomial struct {
	Constant *int
	Factors  *[]Factor
}

func NewConstantMonomial(v int) Monomial {
	return Monomial{
		Constant: &v,
		Factors:  nil,
	}
}

func NewProductMonomial(factors []Factor) Monomial {
	return Monomial{
		Constant: nil,
		Factors:  &factors,
	}
}

type Factor struct {
	Variable    string
	Coefficient int
	Power       int
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
