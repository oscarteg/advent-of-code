package day4

import (
	"testing"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/stretchr/testify/require"
)

const invalidPassports = `eyr:1972 cid:100
hcl:#18171d ecl:amb hgt:170 pid:186cm iyr:2018 byr:1926

iyr:2019
hcl:#602927 eyr:1967 hgt:170cm
ecl:grn pid:012533040 byr:1946

hcl:dab227 iyr:2012
ecl:brn hgt:182cm pid:021572410 eyr:2020 byr:1992 cid:277

hgt:59cm ecl:zzz
eyr:2038 hcl:74454a iyr:2023
pid:3556412378 byr:2007`

const validPassports = `pid:087499704 hgt:74in ecl:grn iyr:2012 eyr:2030 byr:1980
hcl:#623a2f

eyr:2029 ecl:blu cid:129 byr:1989
iyr:2014 pid:896056539 hcl:#a97842 hgt:165cm

hcl:#888785
hgt:164cm byr:2001 iyr:2015 cid:88
pid:545766238 ecl:hzl
eyr:2022

iyr:2010 hgt:158cm hcl:#b6652a ecl:blu byr:1944 eyr:2021 pid:093154719`

func TestB_InvalidPassports(t *testing.T) {
	input := challenge.FromLiteral(invalidPassports)

	result, err := B(input)

	require.NoError(t, err)
	require.Equal(t, 0, result)
}

func TestB_ValidPassports(t *testing.T) {
	input := challenge.FromLiteral(validPassports)

	result, err := B(input)

	require.NoError(t, err)
	require.Equal(t, 4, result)
}

func TestB_MalformedFieldReturnsError(t *testing.T) {
	input := challenge.FromLiteral(
		"byr1937 iyr:2017 eyr:2020 hgt:183cm hcl:#fffffd ecl:gry pid:860033327")

	_, err := B(input)

	require.Error(t, err)
}
