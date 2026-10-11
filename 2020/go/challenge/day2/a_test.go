package day2

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(example)
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 2, result)
}

func TestA_SingleLineEdgeCase(t *testing.T) {
	input := challenge.FromLiteral("2-9 c: ccccccccc")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestA_MalformedInput(t *testing.T) {
	input := challenge.FromLiteral("not a policy line")
	_, err := A(input)
	require.Error(t, err)
}

func TestA_EmptyInput(t *testing.T) {
	input := challenge.FromLiteral("")
	_, err := A(input)
	require.Error(t, err)
}
