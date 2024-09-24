package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorOnEmptyInput(t *testing.T) {
	_, err := Parser{}.Parse("")

	require.Error(t, err)
}
