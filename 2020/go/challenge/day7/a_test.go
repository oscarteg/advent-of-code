package day7

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

func TestAPart1(t *testing.T) {
	example := `light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.`

	answer, err := A(challenge.FromLiteral(example))
	require.NoError(t, err)
	require.Equal(t, 4, answer)
}

func TestAParsing(t *testing.T) {
	rules, err := parseRules([]string{"shiny gold bags contain 2 dark red bags."})
	require.NoError(t, err)
	require.Equal(t, 1, len(rules))
	require.Equal(t, "shiny gold", rules[0].BagType)
	require.Equal(t, 1, len(rules[0].Contents))
	require.Equal(t, 2, rules[0].Contents[0].Count)
	require.Equal(t, "dark red", rules[0].Contents[0].Type)
}

func TestANoContents(t *testing.T) {
	rules, err := parseRules([]string{"faded blue bags contain no other bags."})
	require.NoError(t, err)
	require.Equal(t, 0, len(rules[0].Contents))
}

func TestAMalformed(t *testing.T) {
	_, err := parseRules([]string{})
	require.Error(t, err)
}
