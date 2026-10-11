package day15

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	result, err := play([]int{0, 3, 6}, 30_000_000)

	require.NoError(t, err)
	require.Equal(t, 175594, result)
}

func TestBTurnEqualToStartingLengthReturnsLastStartingNumber(t *testing.T) {
	result, err := play([]int{0, 3, 6}, 3)

	require.NoError(t, err)
	require.Equal(t, 6, result)
}

func TestBMalformedInput(t *testing.T) {
	input := challenge.FromLiteral("")

	_, err := B(input)

	require.Error(t, err)
}
