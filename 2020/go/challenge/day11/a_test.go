package day11

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func TestAExample(t *testing.T) {
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
	answer, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 37, answer)
}

func TestASingleSeat(t *testing.T) {
	input := challenge.FromLiteral("L")
	answer, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 1, answer)
}

func TestAFloorOnly(t *testing.T) {
	input := challenge.FromLiteral(".")
	answer, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 0, answer)
}

func TestAInvalidCell(t *testing.T) {
	input := challenge.FromLiteral("X")
	_, err := A(input)
	require.Error(t, err)
}
