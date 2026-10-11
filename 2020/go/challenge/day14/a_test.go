package day14

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestPartA(t *testing.T) {
	input := challenge.FromLiteral(`mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X
mem[8] = 11
mem[7] = 101
mem[8] = 0`)
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, uint64(165), result)
}

func TestPartAMaskApplication(t *testing.T) {
	masked := applyMask(11, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X")
	require.Equal(t, uint64(73), masked)
}

func TestPartAMalformed(t *testing.T) {
	input := challenge.FromLiteral("invalid input")
	result, err := A(input)
	require.NoError(t, err)
	require.Equal(t, uint64(0), result)
}
