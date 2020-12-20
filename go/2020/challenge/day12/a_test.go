package day12

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestA(t *testing.T) {
	input := challenge.FromLiteral("foobar")
	result := a(input)
	require.Equal(t, 42, result)
}
