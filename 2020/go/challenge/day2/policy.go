package day2

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

var lineRE = regexp.MustCompile(`^(\d+)-(\d+) (\w): (\w+)$`)

type policy struct {
	low, high int
	letter    byte
	password  string
}

func parsePolicy(line string) (policy, error) {
	match := lineRE.FindStringSubmatch(line)
	if match == nil {
		return policy{}, fmt.Errorf("line does not match policy format: %q", line)
	}

	low, err := strconv.Atoi(match[1])
	if err != nil {
		return policy{}, fmt.Errorf("low bound is not a number: %w", err)
	}

	high, err := strconv.Atoi(match[2])
	if err != nil {
		return policy{}, fmt.Errorf("high bound is not a number: %w", err)
	}

	return policy{
		low:      low,
		high:     high,
		letter:   match[3][0],
		password: match[4],
	}, nil
}

func parsePolicies(input *challenge.Input) ([]policy, error) {
	var policies []policy
	for line := range input.Lines() {
		p, err := parsePolicy(line)
		if err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}

	if len(policies) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	return policies, nil
}

func (p policy) isValidByCount() bool {
	var occurrences int
	for i := 0; i < len(p.password); i++ {
		if p.password[i] == p.letter {
			occurrences++
		}
	}

	return occurrences >= p.low && occurrences <= p.high
}

func (p policy) isValidByPosition() (bool, error) {
	at := func(pos int) (byte, error) {
		if pos < 1 || pos > len(p.password) {
			return 0, fmt.Errorf("password %q has no character at position %d", p.password, pos)
		}
		return p.password[pos-1], nil
	}

	low, err := at(p.low)
	if err != nil {
		return false, err
	}
	high, err := at(p.high)
	if err != nil {
		return false, err
	}

	return (low == p.letter) != (high == p.letter), nil
}
