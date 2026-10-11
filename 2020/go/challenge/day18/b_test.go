package day18

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB_PerLineExamples(t *testing.T) {
	lines := []string{
		"1 + 2 * 3 + 4 * 5 + 6",
		"1 + (2 * 3) + (4 * (5 + 6))",
		"2 * 3 + (4 * 5)",
		"5 + (8 * 3 + 9 + 3 * 4 * 3)",
		"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))",
		"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2",
	}
	expected := []int{231, 51, 46, 1445, 669060, 23340}

	for i, line := range lines {
		value, err := evalExpression(line, true)
		require.NoError(t, err)
		require.Equal(t, expected[i], value, line)
	}
}

func TestB_Example(t *testing.T) {
	input := challenge.FromLiteral(day18Example)
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 231+51+46+1445+669060+23340, result)
}

func TestB_SingleNumber(t *testing.T) {
	input := challenge.FromLiteral("42")
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 42, result)
}

func TestB_TrailingTokens(t *testing.T) {
	input := challenge.FromLiteral("(1 + 2))")
	_, err := B(input)
	require.Error(t, err)
}

func TestB_EmptyInput(t *testing.T) {
	input := challenge.FromLiteral("")
	_, err := B(input)
	require.Error(t, err)
}
