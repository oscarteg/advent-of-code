package day1

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(day1example)

	result := a(input)

	require.Equal(t, 514579, result)
}
