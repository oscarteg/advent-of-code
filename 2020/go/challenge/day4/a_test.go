package day4

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const day4Example = `ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in`

func TestA(t *testing.T) {
	input := challenge.FromLiteral(day4Example)

	result, err := A(input)

	require.NoError(t, err)
	require.Equal(t, 2, result)
}

func TestA_OnlyCIDMissingIsValid(t *testing.T) {
	input := challenge.FromLiteral(
		"byr:1937 iyr:2017 eyr:2020 hgt:183cm hcl:#fffffd ecl:gry pid:860033327")

	result, err := A(input)

	require.NoError(t, err)
	require.Equal(t, 1, result)
}

func TestA_MalformedFieldReturnsError(t *testing.T) {
	input := challenge.FromLiteral(
		"byr1937 iyr:2017 eyr:2020 hgt:183cm hcl:#fffffd ecl:gry pid:860033327")

	_, err := A(input)

	require.Error(t, err)
}

func TestA_EmptyInputReturnsError(t *testing.T) {
	input := challenge.FromLiteral("")

	_, err := A(input)

	require.Error(t, err)
}
