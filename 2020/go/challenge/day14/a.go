package day14

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (uint64, error) {
	memory := make(map[uint64]uint64)
	var mask string

	maskRegex := regexp.MustCompile(`^mask = ([01X]+)$`)
	memRegex := regexp.MustCompile(`^mem\[(\d+)\] = (\d+)$`)

	for line := range input.Lines() {
		line = strings.TrimSpace(line)
		if maskMatches := maskRegex.FindStringSubmatch(line); maskMatches != nil {
			mask = maskMatches[1]
		} else if memMatches := memRegex.FindStringSubmatch(line); memMatches != nil {
			addr, _ := strconv.ParseUint(memMatches[1], 10, 64)
			value, _ := strconv.ParseUint(memMatches[2], 10, 64)
			memory[addr] = applyMask(value, mask)
		}
	}

	var sum uint64
	for _, v := range memory {
		sum += v
	}
	return sum, nil
}

func applyMask(value uint64, mask string) uint64 {
	result := value
	for i := 0; i < len(mask); i++ {
		bitPos := len(mask) - 1 - i
		switch mask[i] {
		case '0':
			result &= ^(1 << uint(bitPos))
		case '1':
			result |= 1 << uint(bitPos)
		}
	}
	return result
}
