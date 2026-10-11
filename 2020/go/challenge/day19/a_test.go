package day19

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const example = `0: 4 1 5
1: 2 3 | 3 2
2: 4 4 | 5 5
3: 4 5 | 5 4
4: "a"
5: "b"

ababbb
bababa
abbbab
aaabbb
aaaabbb
`

// The larger official example; its rules 8 and 11 loop in part 2.
const loopExample = `42: 9 14 | 10 1
9: 14 27 | 1 26
10: 23 14 | 28 1
1: "a"
11: 42 31
5: 1 14 | 15 1
19: 14 1 | 14 14
12: 24 14 | 19 1
16: 15 1 | 14 14
31: 14 17 | 1 13
6: 14 14 | 1 14
2: 1 24 | 14 4
0: 8 11
13: 14 3 | 1 12
15: 1 | 14
17: 14 2 | 1 7
23: 25 1 | 22 14
28: 16 1
4: 1 1
20: 14 14 | 1 15
3: 5 14 | 16 1
27: 1 6 | 14 18
14: "b"
21: 14 1 | 1 14
25: 1 1 | 1 14
22: 14 14
8: 42
26: 14 22 | 1 20
18: 15 15
7: 14 5 | 1 21
24: 14 1

abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa
bbabbbbaabaabba
babbbbaabbbbbabbbbbbaabaaabaaa
aaabbbbbbaaaabaababaabababbabaaabbababababaaa
bbbbbbbaaaabbbbaaabbabaaa
bbbababbbbaaaaaaaabbababaaababaabab
ababaaaaaabaaab
ababaaaaabbbaba
baabbaaaabbaaaababbaababb
abbbbabbbbaaaababbbbbbaaaababb
aaaaabbaabaaaaababaa
aaaabbaaaabbaaa
aaaabbaabbaaaaaaabbbabbbaaabbaabaab
babaaabbbaaabaababbaabababaaab
aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba
`

func TestA(t *testing.T) {
	answer, err := A(challenge.FromLiteral(example))

	require.NoError(t, err)
	require.Equal(t, 2, answer)
}

func TestALoopExample(t *testing.T) {
	answer, err := A(challenge.FromLiteral(loopExample))

	require.NoError(t, err)
	require.Equal(t, 3, answer)
}

func TestAAlternativesOfDifferentLengthsBothCount(t *testing.T) {
	answer, err := A(challenge.FromLiteral("0: 1 | 1 1\n1: \"a\"\n\na\naa\naaa\n"))

	require.NoError(t, err)
	require.Equal(t, 2, answer)
}

func TestALeftoverCharactersDoNotMatch(t *testing.T) {
	answer, err := A(challenge.FromLiteral("0: 1 1\n1: \"a\"\n\naa\naaa\n"))

	require.NoError(t, err)
	require.Equal(t, 1, answer)
}

func TestAMalformedInput(t *testing.T) {
	for _, input := range []string{
		"",
		"0: 4 1 5\n\nab",
		"0: \"ab\"\n\nab",
		"0 1\n1: \"a\"\n\na",
		"x: \"a\"\n\na",
		"0: 1\n1: \"a\"\n1: \"b\"\n\na",
		"1: \"a\"\n\na",
		"0: \"a\"\n",
		"0: | 1\n1: \"a\"\n\na",
	} {
		_, err := A(challenge.FromLiteral(input))

		require.Error(t, err, "expected an error for %q", input)
	}
}
