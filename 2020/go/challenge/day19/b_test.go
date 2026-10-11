package day19

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	answer, err := B(challenge.FromLiteral(loopExample))

	require.NoError(t, err)
	require.Equal(t, 12, answer)
}

func TestBLeavesARuleSetWithoutLoopsAlone(t *testing.T) {
	answer, err := B(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 2, answer)
}

func TestBLoopRulesWithoutTheRulesTheyNeed(t *testing.T) {
	_, err := B(challenge.FromLiteral("0: 8 11\n8: 9\n11: 9 9\n9: \"a\"\n\naaa\n"))

	require.Error(t, err)
}

func TestBMalformedInput(t *testing.T) {
	_, err := B(challenge.FromLiteral("0: \"ab\"\n\nab"))

	require.Error(t, err)
}
