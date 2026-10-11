package day14

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestPartB(t *testing.T) {
	input := challenge.FromLiteral(`mask = 000000000000000000000000000000X1001X
mem[42] = 100
mask = 00000000000000000000000000000000X0XX
mem[26] = 1`)
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, uint64(208), result)
}

func TestPartBAddressDecoding(t *testing.T) {
	addresses, err := decodeAddresses(27, "000000000000000000000000000000X1101X")
	require.NoError(t, err)
	require.Len(t, addresses, 4)
	require.Contains(t, addresses, uint64(26))
	require.Contains(t, addresses, uint64(27))
	require.Contains(t, addresses, uint64(58))
	require.Contains(t, addresses, uint64(59))
}

func TestPartBMalformed(t *testing.T) {
	input := challenge.FromLiteral("invalid input")
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, uint64(0), result)
}

func TestPartBTooManyFloatingBits(t *testing.T) {
	mask := ""
	for i := 0; i < 25; i++ {
		mask += "X"
	}
	_, err := decodeAddresses(0, mask)
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds maximum")
}
