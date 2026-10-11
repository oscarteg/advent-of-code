package day25

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

const (
	subjectNumber = 7
	modulus       = 20201227
)

func A(input *challenge.Input) (int, error) {
	card, door, err := parsePublicKeys(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	loopSize, err := findLoopSize(card)
	if err != nil {
		return 0, err
	}

	return transform(door, loopSize), nil
}

func parsePublicKeys(lines []string) (card, door int, err error) {
	keys := make([]int, 0, 2)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			return 0, 0, fmt.Errorf("public key %q is not a number: %w", line, err)
		}
		keys = append(keys, value)
	}

	if len(keys) != 2 {
		return 0, 0, fmt.Errorf("expected exactly two public keys, got %d", len(keys))
	}

	return keys[0], keys[1], nil
}

// findLoopSize brute-forces the number of transform iterations needed to
// reach publicKey, since the modulus is small enough to make this fast.
func findLoopSize(publicKey int) (int, error) {
	value := 1
	for loopSize := 1; loopSize <= modulus; loopSize++ {
		value = (value * subjectNumber) % modulus
		if value == publicKey {
			return loopSize, nil
		}
	}
	return 0, fmt.Errorf("could not find a loop size for public key %d", publicKey)
}

func transform(subject, loopSize int) int {
	value := 1
	for i := 0; i < loopSize; i++ {
		value = (value * subject) % modulus
	}
	return value
}
