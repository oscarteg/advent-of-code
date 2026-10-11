package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidReport = errors.New("invalid report")

type history []int64

func parseHistory(line string) (history, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("%w: history has no values", errInvalidReport)
	}
	values := make(history, len(fields))
	for i, field := range fields {
		value, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("history value: %w", err)
		}
		values[i] = value
	}
	return values, nil
}

func parseReport(input string) ([]history, error) {
	var histories []history
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		h, err := parseHistory(line)
		if err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}
	if len(histories) == 0 {
		return nil, fmt.Errorf("%w: report has no histories", errInvalidReport)
	}
	return histories, nil
}

func differences(values history) history {
	diffs := make(history, len(values)-1)
	for i := 1; i < len(values); i++ {
		diffs[i-1] = values[i] - values[i-1]
	}
	return diffs
}

func allZero(values history) bool {
	for _, v := range values {
		if v != 0 {
			return false
		}
	}
	return true
}

func (h history) extrapolateNext() int64 {
	if allZero(h) {
		return 0
	}
	return h[len(h)-1] + differences(h).extrapolateNext()
}

func (h history) extrapolatePrevious() int64 {
	if allZero(h) {
		return 0
	}
	return h[0] - differences(h).extrapolatePrevious()
}

func solve(input string) (int64, int64, error) {
	histories, err := parseReport(input)
	if err != nil {
		return 0, 0, err
	}
	var next, previous int64
	for _, h := range histories {
		next += h.extrapolateNext()
		previous += h.extrapolatePrevious()
	}
	return next, previous, nil
}
