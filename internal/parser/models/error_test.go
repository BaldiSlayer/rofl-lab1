package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapping(t *testing.T) {
	err := func () error {
		return NewParseError("inner error", "вложенная ошибка")
	}()

	err = Wrap("outer error", "внешняя ошибка", err)

	var parseError *ParseError
	assert.ErrorAs(t, err, &parseError)
	assert.Equal(t, "внешняя ошибка: вложенная ошибка", parseError.LlmMessage)
	assert.Equal(t, "outer error: inner error", err.Error())
}
