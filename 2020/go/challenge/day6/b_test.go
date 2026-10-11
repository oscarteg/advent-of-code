package day6

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(example)

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 6, result)
}

func TestBSinglePersonGroup(t *testing.T) {
	input := challenge.FromLiteral("abc")

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 3, result)
}

func TestBEmptyGroupReturnsError(t *testing.T) {
	input := challenge.FromLiteral("abc\n\n\ndef")

	_, err := B(input)
	require.Error(t, err)
}
