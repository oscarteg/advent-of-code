package day2

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(example)
	result := a(input)
	require.Equal(t, 2, result)
}
