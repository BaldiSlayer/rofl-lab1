package models

import "fmt"

type ParseError struct {
	LlmMessage string
	Message    string
}

func (child *ParseError) Wrap(parent *ParseError) *ParseError {
	return &ParseError{
		LlmMessage: fmt.Sprintf("%s: %s", parent.LlmMessage, child.LlmMessage),
		Message:    fmt.Sprintf("%s: %s", parent.Message, child.Message),
	}
}

func (e *ParseError) Error() string {
	return e.Message
}
