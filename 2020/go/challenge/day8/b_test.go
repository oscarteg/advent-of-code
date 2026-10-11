package day8

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	result, err := B(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 8, result)
}

func TestB_NoFixPossibleReturnsError(t *testing.T) {
	_, err := B(challenge.FromLiteral("acc +1\nacc +2\nacc +3"))

	require.Error(t, err)
}

func TestB_MalformedArgumentIsRejected(t *testing.T) {
	_, err := B(challenge.FromLiteral("acc +x"))

	require.Error(t, err)
}
