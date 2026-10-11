package day13

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = "939\n7,13,x,x,59,x,31,19"

func TestA(t *testing.T) {
	input := challenge.FromLiteral(example)
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 295, result)
}

func TestASingleBusWaitsZeroOrFullPeriod(t *testing.T) {
	result, err := A(challenge.FromLiteral("0\n5"))
	require.NoError(t, err)
	require.Equal(t, 0, result)

	result, err = A(challenge.FromLiteral("3\n5"))
	require.NoError(t, err)
	require.Equal(t, 10, result)
}

func TestAMissingBusLine(t *testing.T) {
	_, err := A(challenge.FromLiteral("939"))
	require.Error(t, err)
}

func TestAInvalidBusID(t *testing.T) {
	_, err := A(challenge.FromLiteral("939\n7,abc,59"))
	require.Error(t, err)
}
