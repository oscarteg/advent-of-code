package day10

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestBSmallExample(t *testing.T) {
	result, err := B(challenge.FromLiteral(smallExample))
	require.NoError(t, err)
	require.Equal(t, 8, result)
}

func TestBLargerExample(t *testing.T) {
	result, err := B(challenge.FromLiteral(largerExample))
	require.NoError(t, err)
	require.Equal(t, 19208, result)
}

func TestBSingleAdapterEdgeCase(t *testing.T) {
	result, err := B(challenge.FromLiteral("1"))
	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestBMalformedInput(t *testing.T) {
	_, err := B(challenge.FromLiteral("not-a-number"))
	require.Error(t, err)
}
