package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidInput = errors.New("invalid input")

const modulus = 16777216

func nextSecret(secret int64) int64 {
	secret = ((secret * 64) ^ secret) % modulus
	secret = ((secret / 32) ^ secret) % modulus
	secret = ((secret * 2048) ^ secret) % modulus
	return secret
}

func parseSeeds(input string) ([]int64, error) {
	var seeds []int64
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		seed, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("seed: %w", err)
		}
		seeds = append(seeds, seed)
	}
	if len(seeds) == 0 {
		return nil, fmt.Errorf("%w: no seeds found", errInvalidInput)
	}
	return seeds, nil
}

type sequenceKey [4]int8

func solve(input string) (part1, part2 int64, err error) {
	seeds, err := parseSeeds(input)
	if err != nil {
		return 0, 0, err
	}

	totalBananasBySequence := make(map[sequenceKey]int64)

	for _, seed := range seeds {
		secret := seed
		var deltas sequenceKey
		previousPrice := int8(secret % 10)
		seenBySequence := make(map[sequenceKey]bool)

		for step := 1; step <= 2000; step++ {
			secret = nextSecret(secret)
			price := int8(secret % 10)
			deltas = sequenceKey{deltas[1], deltas[2], deltas[3], price - previousPrice}
			previousPrice = price

			if step >= 4 && !seenBySequence[deltas] {
				seenBySequence[deltas] = true
				totalBananasBySequence[deltas] += int64(price)
			}
		}
		part1 += secret
	}

	for _, bananas := range totalBananasBySequence {
		if bananas > part2 {
			part2 = bananas
		}
	}
	return part1, part2, nil
}
