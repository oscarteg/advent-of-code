package day2

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(example)
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestB_BothPositionsMatchEdgeCase(t *testing.T) {
	input := challenge.FromLiteral("2-9 c: ccccccccc")
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestB_MalformedInput(t *testing.T) {
	input := challenge.FromLiteral("not a policy line")
	_, err := B(input)
	require.Error(t, err)
}
