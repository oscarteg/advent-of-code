package day13

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestB(t *testing.T) {
	input := challenge.FromLiteral(example)
	result, err := B(input)
	require.NoError(t, err)
	require.Equal(t, 1068781, result)
}

func testBLiteral(t *testing.T, busLine string, want int) {
	t.Helper()
	result, err := B(challenge.FromLiteral("0\n" + busLine))
	require.NoError(t, err)
	require.Equal(t, want, result)
}

func TestBAdditionalExamples(t *testing.T) {
	testBLiteral(t, "17,x,13,19", 3417)
	testBLiteral(t, "67,7,59,61", 754018)
	testBLiteral(t, "67,x,7,59,61", 779210)
	testBLiteral(t, "67,7,x,59,61", 1261476)
	testBLiteral(t, "1789,37,47,1889", 1202161486)
}

func TestBAllXBusLineIsRejected(t *testing.T) {
	_, err := B(challenge.FromLiteral("939\nx,x,x"))
	require.Error(t, err)
}
