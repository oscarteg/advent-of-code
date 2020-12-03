package day2

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(example)
	result := b(input)
	require.Equal(t, 1, result)
}
