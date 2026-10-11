package day1

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(day1example)

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 241861950, result)
}

func TestB_InvalidInput(t *testing.T) {
	input := challenge.FromLiteral("100\n200\ninvalid\n400")
	_, err := B(input)
	require.Error(t, err)
}

func TestB_NoSolution(t *testing.T) {
	input := challenge.FromLiteral("1\n2\n3\n4\n5")
	_, err := B(input)
	require.Error(t, err)
}
