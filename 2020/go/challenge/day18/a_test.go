package day18

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day18Example = `1 + 2 * 3 + 4 * 5 + 6
1 + (2 * 3) + (4 * (5 + 6))
2 * 3 + (4 * 5)
5 + (8 * 3 + 9 + 3 * 4 * 3)
5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))
((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2
`

func TestA_PerLineExamples(t *testing.T) {
	lines := []string{
		"1 + 2 * 3 + 4 * 5 + 6",
		"1 + (2 * 3) + (4 * (5 + 6))",
		"2 * 3 + (4 * 5)",
		"5 + (8 * 3 + 9 + 3 * 4 * 3)",
		"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))",
		"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2",
	}
	expected := []int{71, 51, 26, 437, 12240, 13632}

	for i, line := range lines {
		value, err := evalExpression(line, false)
		require.NoError(t, err)
		require.Equal(t, expected[i], value, line)
	}
}

func TestA_Example(t *testing.T) {
	input := challenge.FromLiteral(day18Example)
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 71+51+26+437+12240+13632, result)
}

func TestA_SingleNumber(t *testing.T) {
	input := challenge.FromLiteral("42")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 42, result)
}

func TestA_InvalidCharacter(t *testing.T) {
	input := challenge.FromLiteral("1 + a")
	_, err := A(input)
	require.Error(t, err)
}

func TestA_UnmatchedParenthesis(t *testing.T) {
	input := challenge.FromLiteral("(1 + 2")
	_, err := A(input)
	require.Error(t, err)
}

func TestA_EmptyInput(t *testing.T) {
	input := challenge.FromLiteral("")
	_, err := A(input)
	require.Error(t, err)
}
