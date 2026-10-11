package day22

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(example)

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 306, result)
}

func TestA_SingleCardDeck(t *testing.T) {
	input := challenge.FromLiteral("Player 1:\n5\n\nPlayer 2:\n3")

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 13, result)
}

func TestA_MalformedInput(t *testing.T) {
	input := challenge.FromLiteral("Player 1:\n9\nabc\n\nPlayer 2:\n5\n4")

	_, err := A(input)
	require.Error(t, err)
}

func TestA_MissingPlayerBlock(t *testing.T) {
	input := challenge.FromLiteral("Player 1:\n9\n2\n6\n3\n1")

	_, err := A(input)
	require.Error(t, err)
}
