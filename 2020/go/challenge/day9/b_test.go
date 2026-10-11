package day9

import (
	"errors"
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	result, err := B(challenge.FromLiteral(day9example), 5)
	require.NoError(t, err)
	require.Equal(t, int64(62), result)
}

func TestBEdgeCase(t *testing.T) {
	result, err := B(challenge.FromLiteral("1\n2\n3\n6"), 2)
	require.NoError(t, err)
	require.Equal(t, int64(4), result)
}

func TestBNoContiguousRun(t *testing.T) {
	_, err := findWeakness([]int64{1, 5, 9}, 4)
	require.True(t, errors.Is(err, errNotFound), "got %v", err)
}

func TestBInvalidInput(t *testing.T) {
	_, err := B(challenge.FromLiteral("1\nabc\n3"), 2)
	require.Error(t, err)
}
