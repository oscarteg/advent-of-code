package main

import (
	"fmt"
	"strings"
)

func snafuDigitValue(digit byte) (int64, error) {
	switch digit {
	case '2':
		return 2, nil
	case '1':
		return 1, nil
	case '0':
		return 0, nil
	case '-':
		return -1, nil
	case '=':
		return -2, nil
	default:
		return 0, fmt.Errorf("invalid SNAFU digit: %q", digit)
	}
}

func parseSnafu(line string) (int64, error) {
	if line == "" {
		return 0, fmt.Errorf("SNAFU number must not be empty")
	}
	var value int64
	for i := 0; i < len(line); i++ {
		digit, err := snafuDigitValue(line[i])
		if err != nil {
			return 0, fmt.Errorf("malformed SNAFU number %q: %w", line, err)
		}
		value = value*5 + digit
	}
	return value, nil
}

func toSnafu(value int64) string {
	if value == 0 {
		return "0"
	}
	var digits []byte
	for value != 0 {
		remainder := ((value % 5) + 5) % 5
		var digit byte
		var carry int64
		switch remainder {
		case 0:
			digit, carry = '0', 0
		case 1:
			digit, carry = '1', 0
		case 2:
			digit, carry = '2', 0
		case 3:
			digit, carry = '=', 1
		case 4:
			digit, carry = '-', 1
		}
		digits = append(digits, digit)
		value = value/5 + carry
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}

func solve(input string) (string, error) {
	var sum int64
	count := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		value, err := parseSnafu(line)
		if err != nil {
			return "", err
		}
		sum += value
		count++
	}
	if count == 0 {
		return "", fmt.Errorf("input must contain at least one SNAFU number")
	}
	return toSnafu(sum), nil
}
