package day14

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

const maxFloatingBits = 20

func B(input *challenge.Input) (uint64, error) {
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
			addresses, err := decodeAddresses(addr, mask)
			if err != nil {
				return 0, err
			}
			for _, a := range addresses {
				memory[a] = value
			}
		}
	}

	var sum uint64
	for _, v := range memory {
		sum += v
	}
	return sum, nil
}

func decodeAddresses(address uint64, mask string) ([]uint64, error) {
	var floatingPos []int
	baseAddr := address

	for i := 0; i < len(mask); i++ {
		bitPos := len(mask) - 1 - i
		switch mask[i] {
		case '0':
		case '1':
			baseAddr |= 1 << uint(bitPos)
		case 'X':
			floatingPos = append(floatingPos, bitPos)
			baseAddr &= ^(1 << uint(bitPos))
		}
	}

	numFloating := len(floatingPos)
	if numFloating > maxFloatingBits {
		return nil, fmt.Errorf("mask has %d floating bits, exceeds maximum of %d", numFloating, maxFloatingBits)
	}

	addresses := make([]uint64, 0, 1<<uint(numFloating))

	for variant := uint64(0); variant < (1 << uint(numFloating)); variant++ {
		a := baseAddr
		for j, pos := range floatingPos {
			if (variant>>uint(j))&1 == 1 {
				a |= 1 << uint(pos)
			}
		}
		addresses = append(addresses, a)
	}

	return addresses, nil
}
