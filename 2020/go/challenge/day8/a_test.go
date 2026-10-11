package day8

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6`

func TestA(t *testing.T) {
	result, err := A(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 5, result)
}

func TestA_TerminatingProgramKeepsItsAccumulator(t *testing.T) {
	result, err := A(challenge.FromLiteral("acc +1\nacc +2\nacc +3"))

	require.NoError(t, err)
	require.Equal(t, 6, result)
}

func TestA_MalformedInstructionIsRejected(t *testing.T) {
	_, err := A(challenge.FromLiteral("foo +1"))

	require.Error(t, err)
}
