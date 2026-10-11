package day3

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(example)
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 336, result)
}

func TestBEmptyInput(t *testing.T) {
	input := challenge.FromLiteral("")
	_, err := B(input)
	require.Error(t, err)
}

func TestBSingleLine(t *testing.T) {
	input := challenge.FromLiteral("#..")
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 1, result)
}
