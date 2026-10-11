package day25

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day25example = `5764801
17807724
`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(day25example)

	result, err := A(input)

	require.NoError(t, err)
	require.Equal(t, 14897079, result)
}

func TestFindLoopSizeMatchesKnownValues(t *testing.T) {
	cardLoopSize, err := findLoopSize(5764801)
	require.NoError(t, err)
	require.Equal(t, 8, cardLoopSize)

	doorLoopSize, err := findLoopSize(17807724)
	require.NoError(t, err)
	require.Equal(t, 11, doorLoopSize)
}

func TestDoorLoopSizeGivesSameKey(t *testing.T) {
	card, door, err := parsePublicKeys([]string{"5764801", "17807724"})
	require.NoError(t, err)

	loopSize, err := findLoopSize(door)
	require.NoError(t, err)

	require.Equal(t, 14897079, transform(card, loopSize))
}

func TestAMissingPublicKeyIsAnError(t *testing.T) {
	input := challenge.FromLiteral("5764801\n")

	_, err := A(input)

	require.Error(t, err)
}

func TestANonNumericPublicKeyIsAnError(t *testing.T) {
	input := challenge.FromLiteral("not-a-number\n17807724\n")

	_, err := A(input)

	require.Error(t, err)
}

func TestAExtraLineIsAnError(t *testing.T) {
	input := challenge.FromLiteral("5764801\n17807724\n123\n")

	_, err := A(input)

	require.Error(t, err)
}
