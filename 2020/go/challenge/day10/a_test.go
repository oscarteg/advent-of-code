package day10

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const smallExample = `16
10
15
5
1
11
7
19
6
12
4`

const largerExample = `28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3`

func TestASmallExample(t *testing.T) {
	result, err := A(challenge.FromLiteral(smallExample))
	require.NoError(t, err)
	require.Equal(t, 35, result)
}

func TestALargerExample(t *testing.T) {
	result, err := A(challenge.FromLiteral(largerExample))
	require.NoError(t, err)
	require.Equal(t, 220, result)
}

func TestASingleAdapterEdgeCase(t *testing.T) {
	// Outlet (0) -> adapter (1) -> device (4): one 1-jolt gap, one 3-jolt gap.
	result, err := A(challenge.FromLiteral("1"))
	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestAMalformedInput(t *testing.T) {
	_, err := A(challenge.FromLiteral("not-a-number"))
	require.Error(t, err)
}

func TestAEmptyInput(t *testing.T) {
	_, err := A(challenge.FromLiteral(""))
	require.Error(t, err)
}
