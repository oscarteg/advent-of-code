package day5

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestBMissingOddSeat(t *testing.T) {
	input := `FBFBBFFRLR
FBFBBFFRRL
FBFBBFBLLL
FBFBBFBLLR`

	result, err := B(challenge.FromLiteral(input))
	require.NoError(t, err)
	require.Equal(t, 359, result)
}

func TestBMissingEvenSeat(t *testing.T) {
	input := `FFFBBBFRRR
FFFBBBBLLL
FFFBBBBLRL
FFFBBBBLRR`

	result, err := B(challenge.FromLiteral(input))
	require.NoError(t, err)
	require.Equal(t, 121, result)
}

func TestBNoMissingSeat(t *testing.T) {
	input := `FBFBBFFRLR
FBFBBFFRRR
FBFBBFFRRL
FBFBBFFRRR`

	_, err := B(challenge.FromLiteral(input))
	require.Error(t, err)
}

func TestBMissingFromBoarding(t *testing.T) {
	input := `FBFBBFFRLL
FBFBBFFRLR
FBFBBFFRRR
FBFBBFBLLL`

	result, err := B(challenge.FromLiteral(input))
	require.NoError(t, err)
	require.Equal(t, 358, result)
}
