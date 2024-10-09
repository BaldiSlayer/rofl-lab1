package trsparser

import (
	"errors"

	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
)

func (e *ParseError) Error() string {
	return e.Summary
}

func toParseError(err error) error {
	var parseError *models.ParseError
	if errors.As(err, &parseError) {
		return &ParseError{
			LlmMessage: parseError.LlmMessage,
			Summary:    parseError.Message,
		}
	}
	return &ParseError{
		LlmMessage: err.Error(),
		Summary:    err.Error(),
	}
}
