package day21

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(day21Example)

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, "mxmxvkd,sqjhc,fvjkl", result)
}

func TestB_SingleFoodSingleIngredientIsAnEdgeCase(t *testing.T) {
	input := challenge.FromLiteral("a (contains dairy)")

	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, "a", result)
}

func TestB_AmbiguousAssignmentIsError(t *testing.T) {
	input := challenge.FromLiteral("a b (contains dairy, fish)")

	_, err := B(input)
	require.Error(t, err)
}

func TestB_EmptyInputIsError(t *testing.T) {
	input := challenge.FromLiteral("")

	_, err := B(input)
	require.Error(t, err)
}
