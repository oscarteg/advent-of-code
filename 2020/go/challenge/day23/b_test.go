package day23

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	result, err := B(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 149245887792, result)
}

func TestB_PaddingCupsContinueTheCircle(t *testing.T) {
	cups, err := parseCups(challenge.FromLiteral(example))
	require.NoError(t, err)

	next, err := play(cups, 12, 0)
	require.NoError(t, err)

	circle := make([]int, 0, 12)
	cup := cups[0]
	for range 12 {
		circle = append(circle, cup)
		cup = next[cup]
	}

	require.Equal(t, []int{3, 8, 9, 1, 2, 5, 4, 6, 7, 10, 11, 12}, circle)
	require.Equal(t, cups[0], cup, "the circle closes onto the first cup")
}

func TestB_MalformedInputIsRejected(t *testing.T) {
	_, err := B(challenge.FromLiteral("389125467\n213456789"))

	require.Error(t, err)
}
