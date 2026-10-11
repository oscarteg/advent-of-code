package day16

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

// The second example of the puzzle: it has no departure fields, but its field order is fully
// determined.
const fieldOrderExample = `class: 0-1 or 4-19
row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9`

func TestB_SecondExampleResolvesEveryField(t *testing.T) {
	parsed, err := parseNotes(challenge.FromLiteral(fieldOrderExample).LinesSlice())
	require.NoError(t, err)

	order, err := parsed.fieldOrder()
	require.NoError(t, err)

	fields := make(map[string]int, len(order))
	for column, name := range order {
		fields[name] = parsed.ticket[column]
	}
	require.Equal(t, map[string]int{"class": 12, "row": 11, "seat": 13}, fields)
}

func TestB_DepartureFieldsAreMultiplied(t *testing.T) {
	input := `departure gate: 1-3
departure track: 1-6
seat: 1-9

your ticket:
2,5,7

nearby tickets:
1,4,9
3,6,8`

	result, err := B(challenge.FromLiteral(input))

	require.NoError(t, err)
	require.Equal(t, 10, result)
}

func TestB_ExampleWithoutDepartureFieldsIsAnEmptyProduct(t *testing.T) {
	result, err := B(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestB_AmbiguousFieldOrderIsRejected(t *testing.T) {
	input := `first: 1-2
second: 1-2

your ticket:
1,2

nearby tickets:
1,2`

	_, err := B(challenge.FromLiteral(input))

	require.Error(t, err)
}
