package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Packet interface{}

func parsePacket(s string) (Packet, error) {
	var p Packet
	err := json.Unmarshal([]byte(s), &p)
	return p, err
}

func compare(a, b Packet) int {
	switch av := a.(type) {
	case float64:
		switch bv := b.(type) {
		case float64:
			if av < bv {
				return -1
			} else if av > bv {
				return 1
			}
			return 0
		case []interface{}:
			aList := []interface{}{av}
			return compareSlices(aList, bv)
		}
	case []interface{}:
		switch bv := b.(type) {
		case float64:
			bList := []interface{}{bv}
			return compareSlices(av, bList)
		case []interface{}:
			return compareSlices(av, bv)
		}
	}
	return 0
}

func compareSlices(a, b []interface{}) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		cmp := compare(a[i], b[i])
		if cmp != 0 {
			return cmp
		}
	}

	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return 0
}

func solve(input string) (int, int, error) {
	part1, err := solvePart1(input)
	if err != nil {
		return 0, 0, err
	}

	part2, err := solvePart2(input)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

func solvePart1(input string) (int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	sum := 0
	pairIdx := 1

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		packet1, err := parsePacket(line)
		if err != nil {
			return 0, fmt.Errorf("failed to parse packet: %w", err)
		}

		if i+1 >= len(lines) {
			return 0, fmt.Errorf("incomplete pair at index %d", pairIdx)
		}

		i++
		line2 := strings.TrimSpace(lines[i])
		if line2 == "" {
			i--
			continue
		}

		packet2, err := parsePacket(line2)
		if err != nil {
			return 0, fmt.Errorf("failed to parse packet: %w", err)
		}

		if compare(packet1, packet2) < 0 {
			sum += pairIdx
		}
		pairIdx++
	}

	return sum, nil
}

func solvePart2(input string) (int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var packets []Packet

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		packet, err := parsePacket(line)
		if err != nil {
			return 0, fmt.Errorf("failed to parse packet: %w", err)
		}
		packets = append(packets, packet)
	}

	divider1 := []interface{}{[]interface{}{float64(2)}}
	divider2 := []interface{}{[]interface{}{float64(6)}}

	packets = append(packets, divider1, divider2)

	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) < 0
	})

	var pos1, pos2 int
	for i, p := range packets {
		if isDiv1(p) {
			pos1 = i + 1
		}
		if isDiv2(p) {
			pos2 = i + 1
		}
	}

	return pos1 * pos2, nil
}

func isDiv1(p Packet) bool {
	list, ok := p.([]interface{})
	if !ok || len(list) != 1 {
		return false
	}
	inner, ok := list[0].([]interface{})
	if !ok || len(inner) != 1 {
		return false
	}
	val, ok := inner[0].(float64)
	return ok && val == 2
}

func isDiv2(p Packet) bool {
	list, ok := p.([]interface{})
	if !ok || len(list) != 1 {
		return false
	}
	inner, ok := list[0].([]interface{})
	if !ok || len(inner) != 1 {
		return false
	}
	val, ok := inner[0].(float64)
	return ok && val == 6
}
