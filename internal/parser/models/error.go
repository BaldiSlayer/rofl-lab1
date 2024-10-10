package models

import (
	"fmt"
)

type ParseError struct {
	LlmMessage string
	Message    string
	Err        error
}

func NewParseError(message, llmMessage string) error {
	return &ParseError{
		LlmMessage: llmMessage,
		Message:    message,
		Err:        nil,
	}
}

func (e *ParseError) Error() string {
	return e.Message
}

func (err ParseError) Unwrap() error {
	return err.Err
}

func Wrap(message, llmMessage string, child error) error {
	if parseError, ok := child.(*ParseError); ok {
		return &ParseError{
			LlmMessage: fmt.Sprintf("%s: %s", llmMessage, parseError.LlmMessage),
			Message:    fmt.Sprintf("%s: %s", message, parseError.Message),
			Err:        child,
		}
	}
	return &ParseError{
		LlmMessage: llmMessage,
		Message:    fmt.Sprintf("%s: %s", message, child.Error()),
		Err:        child,
	}
}
