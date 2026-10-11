package day22

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(example)

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 291, result)
}

func TestB_InfiniteGameTerminates(t *testing.T) {
	// Plain Combat cycles forever on this deck; Recursive Combat's seen-states
	// rule must still terminate it.
	input := challenge.FromLiteral("Player 1:\n43\n19\n\nPlayer 2:\n2\n29\n14")

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 105, result)
}

func TestB_MalformedInput(t *testing.T) {
	input := challenge.FromLiteral("Player 1:\n\nPlayer 2:\n5\n4")

	_, err := B(input)
	require.Error(t, err)
}
