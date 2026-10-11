package day1

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day1example = `1721
979
366
299
675
1456
`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(day1example)

	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 514579, result)
}

func TestA_SimplePair(t *testing.T) {
	input := challenge.FromLiteral("1000\n1020\n500\n600\n700")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, 1000*1020, result)
}

func TestA_InvalidInput(t *testing.T) {
	input := challenge.FromLiteral("123\nabc\n456")
	_, err := A(input)
	require.Error(t, err)
}

func TestA_NoSolution(t *testing.T) {
	input := challenge.FromLiteral("1\n2\n3\n4\n5")
	_, err := A(input)
	require.Error(t, err)
}
