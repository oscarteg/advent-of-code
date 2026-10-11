package day16

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `class: 1-3 or 5-7
row: 6-11 or 33-44
seat: 13-40 or 45-50

your ticket:
7,1,14

nearby tickets:
7,3,47
40,4,50
55,2,20
38,6,12`

func TestA(t *testing.T) {
	result, err := A(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 71, result)
}

func TestA_TicketsWithinTheRangesHaveNoErrorRate(t *testing.T) {
	result, err := A(challenge.FromLiteral(fieldOrderExample))

	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestA_MalformedInputIsRejected(t *testing.T) {
	cases := map[string]string{
		"rule without a colon":   "class 1-3\n\nyour ticket:\n1\n\nnearby tickets:\n2",
		"inverted range":         "class: 5-1\n\nyour ticket:\n1\n\nnearby tickets:\n2",
		"missing section":        "class: 1-3\n\nyour ticket:\n1",
		"wrong header":           "class: 1-3\n\nthe ticket:\n1\n\nnearby tickets:\n2",
		"non-numeric value":      "class: 1-3\n\nyour ticket:\n1\n\nnearby tickets:\nx",
		"wrong number of values": "class: 1-3\n\nyour ticket:\n1,2\n\nnearby tickets:\n2,3",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := A(challenge.FromLiteral(input))

			require.Error(t, err)
		})
	}
}
