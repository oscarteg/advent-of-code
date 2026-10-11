package day17

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(day17example)
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 848, result)
}

func TestBSingleCube(t *testing.T) {
	input := challenge.FromLiteral("#")
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestBEmpty(t *testing.T) {
	input := challenge.FromLiteral("")
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestBMalformed(t *testing.T) {
	_, err := B(challenge.FromLiteral(".#.\n.x.\n"))
	require.Error(t, err)
}
