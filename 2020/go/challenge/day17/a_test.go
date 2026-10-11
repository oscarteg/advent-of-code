package day17

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day17example = `.#.
..#
###
`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(day17example)
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 112, result)
}

func TestASingleCube(t *testing.T) {
	input := challenge.FromLiteral("#")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestAEmpty(t *testing.T) {
	input := challenge.FromLiteral("")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestAMalformed(t *testing.T) {
	_, err := A(challenge.FromLiteral(".#.\n.x.\n"))
	require.Error(t, err)
}
