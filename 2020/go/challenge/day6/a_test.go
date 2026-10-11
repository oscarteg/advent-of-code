package day6

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `abc

a
b
c

ab
ac

a
a
a
a

b`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(example)

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 11, result)
}

func TestASinglePersonGroup(t *testing.T) {
	input := challenge.FromLiteral("abc")

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 3, result)
}

func TestAMalformedAnswersReturnsError(t *testing.T) {
	input := challenge.FromLiteral("a1c")

	_, err := A(input)
	require.Error(t, err)
}
