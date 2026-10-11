package day23

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = "389125467"

func TestA(t *testing.T) {
	result, err := A(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, "67384529", result)
}

func TestA_TenMovesMatchTheWalkthrough(t *testing.T) {
	cups, err := parseCups(challenge.FromLiteral(example))
	require.NoError(t, err)

	next, err := play(cups, len(cups), 10)

	require.NoError(t, err)
	require.Equal(t, "92658374", labelsAfterOne(next))
}

func TestA_OneMoveSplicesThreeCupsBehindTheDestination(t *testing.T) {
	cups, err := parseCups(challenge.FromLiteral(example))
	require.NoError(t, err)

	next, err := play(cups, len(cups), 1)

	require.NoError(t, err)
	// After the first move the circle reads 3 2 8 9 1 5 4 6 7.
	require.Equal(t, "54673289", labelsAfterOne(next))
}

func TestA_ZeroMovesKeepASmallCircleInOrder(t *testing.T) {
	cups, err := parseCups(challenge.FromLiteral("32415"))
	require.NoError(t, err)

	next, err := play(cups, len(cups), 0)

	require.NoError(t, err)
	require.Equal(t, "5324", labelsAfterOne(next))
}

func TestA_MalformedInputIsRejected(t *testing.T) {
	for name, input := range map[string]string{
		"empty":              "",
		"blank":              "   ",
		"non digit":          "38912546x",
		"zero label":         "389125460",
		"duplicate label":    "388125467",
		"label out of range": "98765",
		"too few cups":       "3241",
		"two lines":          "389125467\n213456789",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := A(challenge.FromLiteral(input))

			require.Error(t, err)
		})
	}
}
