package day21

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day21Example = `mxmxvkd kfcds sqjhc nhms (contains dairy, fish)
trh fvjkl sbzzf mxmxvkd (contains dairy)
sqjhc fvjkl (contains soy)
sqjhc mxmxvkd sbzzf (contains fish)
`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(day21Example)

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 5, result)
}

func TestA_SingleFoodSingleIngredientIsAnEdgeCase(t *testing.T) {
	input := challenge.FromLiteral("a (contains dairy)")

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestA_MissingContainsSectionIsError(t *testing.T) {
	input := challenge.FromLiteral("mxmxvkd kfcds sqjhc nhms")

	_, err := A(input)
	require.Error(t, err)
}

func TestA_UnclosedAllergenListIsError(t *testing.T) {
	input := challenge.FromLiteral("mxmxvkd kfcds sqjhc nhms (contains dairy, fish")

	_, err := A(input)
	require.Error(t, err)
}

func TestA_EmptyInputIsError(t *testing.T) {
	input := challenge.FromLiteral("")

	_, err := A(input)
	require.Error(t, err)
}
