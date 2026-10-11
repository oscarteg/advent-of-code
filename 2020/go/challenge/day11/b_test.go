package day11

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func TestBExample(t *testing.T) {
	input := challenge.FromLiteral(`L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL`)
	answer, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 26, answer)
}

func TestBSingleSeat(t *testing.T) {
	input := challenge.FromLiteral("L")
	answer, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 1, answer)
}

func TestBFloorOnly(t *testing.T) {
	input := challenge.FromLiteral(".")
	answer, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 0, answer)
}

func TestBInvalidCell(t *testing.T) {
	input := challenge.FromLiteral("X")
	_, err := B(input)
	require.Error(t, err)
}
