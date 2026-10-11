package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidStone = errors.New("invalid stone")

func parseStones(input string) (map[int64]int64, error) {
	stones := make(map[int64]int64)
	for _, field := range strings.Fields(input) {
		value, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errInvalidStone, err)
		}
		if value < 0 {
			return nil, fmt.Errorf("%w: negative value %d", errInvalidStone, value)
		}
		stones[value]++
	}
	if len(stones) == 0 {
		return nil, fmt.Errorf("%w: no stones found", errInvalidStone)
	}
	return stones, nil
}

func digitCount(value int64) int {
	digits := 1
	for value >= 10 {
		value /= 10
		digits++
	}
	return digits
}

func blink(stones map[int64]int64) map[int64]int64 {
	next := make(map[int64]int64, len(stones)*2)
	for value, count := range stones {
		switch {
		case value == 0:
			next[1] += count
		case digitCount(value)%2 == 0:
			divisor := int64(1)
			for i := 0; i < digitCount(value)/2; i++ {
				divisor *= 10
			}
			next[value/divisor] += count
			next[value%divisor] += count
		default:
			next[value*2024] += count
		}
	}
	return next
}

func totalStones(stones map[int64]int64) int64 {
	var total int64
	for _, count := range stones {
		total += count
	}
	return total
}

func solveWith(input string, blinks1, blinks2 int) (part1, part2 int64, err error) {
	stones, err := parseStones(input)
	if err != nil {
		return 0, 0, err
	}

	totalBlinks := blinks1
	if blinks2 > totalBlinks {
		totalBlinks = blinks2
	}
	for blinkNumber := 1; blinkNumber <= totalBlinks; blinkNumber++ {
		stones = blink(stones)
		if blinkNumber == blinks1 {
			part1 = totalStones(stones)
		}
		if blinkNumber == blinks2 {
			part2 = totalStones(stones)
		}
	}
	return part1, part2, nil
}

func solve(input string) (part1, part2 int64, err error) {
	return solveWith(input, 25, 75)
}
