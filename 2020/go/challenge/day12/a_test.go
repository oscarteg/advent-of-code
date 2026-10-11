package day12

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day12Example = `F10
N3
F7
R90
F11
`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(day12Example)

	result, err := A(input)

	require.NoError(t, err)
	require.Equal(t, 25, result)
}

func TestA_FullRotationReturnsToStartHeading(t *testing.T) {
	input := challenge.FromLiteral("R90\nR90\nR90\nR90\nF5\n")

	result, err := A(input)

	require.NoError(t, err)
	require.Equal(t, 5, result)
}

func TestA_MalformedInput(t *testing.T) {
	input := challenge.FromLiteral("Q10\n")

	_, err := A(input)

	require.Error(t, err)
}
