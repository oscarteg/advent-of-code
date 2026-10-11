package day9

import (
	"errors"
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day9example = `35
20
15
25
47
40
62
55
65
95
102
117
150
182
127
219
299
277
309
576
`

func TestA(t *testing.T) {
	result, err := A(challenge.FromLiteral(day9example), 5)
	require.NoError(t, err)
	require.Equal(t, int64(127), result)
}

func TestAEdgeCase(t *testing.T) {
	result, err := A(challenge.FromLiteral("1\n2\n3\n6"), 2)
	require.NoError(t, err)
	require.Equal(t, int64(6), result)
}

func TestAEveryNumberValid(t *testing.T) {
	_, err := A(challenge.FromLiteral("1\n2\n3"), 2)
	require.True(t, errors.Is(err, errNotFound), "got %v", err)
}

func TestAPreambleLongerThanInput(t *testing.T) {
	_, err := A(challenge.FromLiteral(day9example), 25)
	require.Error(t, err)
}

func TestAInvalidInput(t *testing.T) {
	_, err := A(challenge.FromLiteral("1\nabc\n3"), 2)
	require.Error(t, err)
}
