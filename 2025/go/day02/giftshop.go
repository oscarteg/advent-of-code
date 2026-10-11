package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidGiftShop = errors.New("invalid gift shop input")

// idRange is an inclusive range of IDs, parsed from "lo-hi".
type idRange struct {
	lo, hi uint64
}

func parseIDRange(token string) (idRange, error) {
	token = strings.TrimSpace(token)
	lo, hi, found := strings.Cut(token, "-")
	if !found {
		return idRange{}, fmt.Errorf("%w: missing '-' in %q", errInvalidGiftShop, token)
	}
	loValue, err := strconv.ParseUint(lo, 10, 64)
	if err != nil {
		return idRange{}, fmt.Errorf("%w: %w", errInvalidGiftShop, err)
	}
	hiValue, err := strconv.ParseUint(hi, 10, 64)
	if err != nil {
		return idRange{}, fmt.Errorf("%w: %w", errInvalidGiftShop, err)
	}
	if loValue > hiValue {
		return idRange{}, fmt.Errorf("%w: range %q is backwards", errInvalidGiftShop, token)
	}
	return idRange{lo: loValue, hi: hiValue}, nil
}

func parseIDRanges(input string) ([]idRange, error) {
	fields := strings.Split(strings.TrimSpace(input), ",")
	ranges := make([]idRange, 0, len(fields))
	for _, field := range fields {
		r, err := parseIDRange(field)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, r)
	}
	return ranges, nil
}

// isRepeatedExactlyTwice reports whether digits is exactly two copies of a
// shorter sequence, e.g. "6464".
func isRepeatedExactlyTwice(digits string) bool {
	if len(digits)%2 != 0 {
		return false
	}
	half := len(digits) / 2
	return digits[:half] == digits[half:]
}

// hasPeriod reports whether period evenly divides len(digits) and the
// prefix of that length repeats across the whole string.
func hasPeriod(digits string, period int) bool {
	for i := period; i < len(digits); i++ {
		if digits[i] != digits[i%period] {
			return false
		}
	}
	return true
}

// isRepeatedAtLeastTwice reports whether digits is some shorter sequence
// repeated two or more times, e.g. "123123123" or "1111111".
func isRepeatedAtLeastTwice(digits string) bool {
	for period := 1; period < len(digits); period++ {
		if len(digits)%period == 0 && hasPeriod(digits, period) {
			return true
		}
	}
	return false
}

func solve(input string) (part1, part2 uint64, err error) {
	ranges, err := parseIDRanges(input)
	if err != nil {
		return 0, 0, err
	}
	for _, r := range ranges {
		for n := r.lo; ; n++ {
			digits := strconv.FormatUint(n, 10)
			if isRepeatedExactlyTwice(digits) {
				part1 += n
			}
			if isRepeatedAtLeastTwice(digits) {
				part2 += n
			}
			if n == r.hi {
				break
			}
		}
	}
	return part1, part2, nil
}
