package trsparser

import (
	"errors"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/trsinterprets"
)

func (e *ParseError) Error() string {
	return e.Summary
}

func toParseError(err error) error {
	var parseError *trsinterprets.ParseError
	if errors.As(err, &parseError) {
		return &ParseError{
			LlmMessage: parseError.LlmMessage(),
			Summary:    parseError.Error(),
		}
	}
	return &ParseError{
		LlmMessage: err.Error(),
		Summary:    err.Error(),
	}
}
