package day5

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestAExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "FBFBBFFRLR",
			input:    "FBFBBFFRLR",
			expected: 357,
		},
		{
			name:     "BFFFBBFRRR",
			input:    "BFFFBBFRRR",
			expected: 567,
		},
		{
			name:     "FFFBBBFRRR",
			input:    "FFFBBBFRRR",
			expected: 119,
		},
		{
			name:     "BBFFBBFRLL",
			input:    "BBFFBBFRLL",
			expected: 820,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := A(challenge.FromLiteral(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAMultiplePasses(t *testing.T) {
	input := `FBFBBFFRLR
BFFFBBFRRR
FFFBBBFRRR
BBFFBBFRLL`

	result, err := A(challenge.FromLiteral(input))
	require.NoError(t, err)
	require.Equal(t, 820, result)
}

func TestAParseSeatMalformed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "too short",
			input: "FBFBBFFRL",
		},
		{
			name:  "invalid row char",
			input: "XBFBBFFRLR",
		},
		{
			name:  "invalid col char",
			input: "FBFBBFFXLR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSeat(tt.input)
			require.Error(t, err)
		})
	}
}
