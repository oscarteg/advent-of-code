package day12

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(day12Example)

	result, err := B(input)

	require.NoError(t, err)
	require.Equal(t, 286, result)
}

func TestB_MalformedValue(t *testing.T) {
	input := challenge.FromLiteral("N10x\n")

	_, err := B(input)

	require.Error(t, err)
}

func TestB_NonRightAngleTurn(t *testing.T) {
	input := challenge.FromLiteral("L45\n")

	_, err := B(input)

	require.Error(t, err)
}
