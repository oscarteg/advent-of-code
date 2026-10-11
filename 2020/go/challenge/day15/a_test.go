package day15

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestA(t *testing.T) {
	input := challenge.FromLiteral("0,3,6")

	result, err := A(input)

	require.NoError(t, err)
	require.Equal(t, 436, result)
}

func TestAMoreExamples(t *testing.T) {
	cases := map[string]int{
		"1,3,2": 1,
		"2,1,3": 10,
		"1,2,3": 27,
		"2,3,1": 78,
		"3,2,1": 438,
		"3,1,2": 1836,
	}

	for literal, expected := range cases {
		input := challenge.FromLiteral(literal)

		result, err := A(input)

		require.NoError(t, err)
		require.Equal(t, expected, result)
	}
}

func TestAMalformedInput(t *testing.T) {
	input := challenge.FromLiteral("1,2,a")

	_, err := A(input)

	require.Error(t, err)
}
