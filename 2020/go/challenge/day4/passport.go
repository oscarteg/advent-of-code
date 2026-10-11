package day4

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

var requiredFields = []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

var validEyeColors = map[string]bool{
	"amb": true, "blu": true, "brn": true, "gry": true, "grn": true, "hzl": true, "oth": true,
}

var (
	heightRE     = regexp.MustCompile(`^(\d+)(cm|in)$`)
	hairColorRE  = regexp.MustCompile(`^#[0-9a-f]{6}$`)
	passportIDRE = regexp.MustCompile(`^\d{9}$`)
)

type passport map[string]string

func parsePassports(input *challenge.Input) ([]passport, error) {
	var blocks []string
	var current []string

	for line := range input.Lines() {
		if line == "" {
			if len(current) > 0 {
				blocks = append(blocks, strings.Join(current, " "))
				current = nil
			}
			continue
		}
		current = append(current, line)
	}
	if len(current) > 0 {
		blocks = append(blocks, strings.Join(current, " "))
	}

	passports := make([]passport, 0, len(blocks))
	for _, block := range blocks {
		p, err := parsePassport(block)
		if err != nil {
			return nil, err
		}
		passports = append(passports, p)
	}

	if len(passports) == 0 {
		return nil, fmt.Errorf("no passports found in input")
	}

	return passports, nil
}

func parsePassport(block string) (passport, error) {
	p := passport{}
	for _, field := range strings.Fields(block) {
		key, value, ok := strings.Cut(field, ":")
		if !ok || key == "" || value == "" {
			return nil, fmt.Errorf("malformed field, missing ':' in %q", field)
		}
		p[key] = value
	}
	return p, nil
}

func (p passport) hasRequiredFields() bool {
	for _, field := range requiredFields {
		if _, ok := p[field]; !ok {
			return false
		}
	}
	return true
}

func (p passport) isValid() bool {
	if !p.hasRequiredFields() {
		return false
	}
	return isValidYear(p["byr"], 1920, 2002) &&
		isValidYear(p["iyr"], 2010, 2020) &&
		isValidYear(p["eyr"], 2020, 2030) &&
		isValidHeight(p["hgt"]) &&
		hairColorRE.MatchString(p["hcl"]) &&
		validEyeColors[p["ecl"]] &&
		passportIDRE.MatchString(p["pid"])
}

func isValidYear(value string, min, max int) bool {
	if len(value) != 4 {
		return false
	}
	year, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return year >= min && year <= max
}

func isValidHeight(value string) bool {
	matches := heightRE.FindStringSubmatch(value)
	if matches == nil {
		return false
	}
	number, err := strconv.Atoi(matches[1])
	if err != nil {
		return false
	}
	switch matches[2] {
	case "cm":
		return number >= 150 && number <= 193
	case "in":
		return number >= 59 && number <= 76
	default:
		return false
	}
}
