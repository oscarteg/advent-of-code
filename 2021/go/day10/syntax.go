package main

import (
	"fmt"
	"sort"
	"strings"
)

type bracket rune

func (b bracket) openChar() rune {
	return rune(b)
}

func (b bracket) closeChar() rune {
	switch b {
	case '(':
		return ')'
	case '[':
		return ']'
	case '{':
		return '}'
	case '<':
		return '>'
	}
	return 0
}

func (b bracket) errorScore() int64 {
	switch b {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	}
	return 0
}

func (b bracket) completionScore() int64 {
	switch b {
	case '(':
		return 1
	case '[':
		return 2
	case '{':
		return 3
	case '<':
		return 4
	}
	return 0
}

func isOpen(ch rune) bool {
	switch ch {
	case '(', '[', '{', '<':
		return true
	}
	return false
}

func isClose(ch rune) bool {
	switch ch {
	case ')', ']', '}', '>':
		return true
	}
	return false
}

func getOpenBracket(ch rune) bracket {
	switch ch {
	case '(', ')':
		return '('
	case '[', ']':
		return '['
	case '{', '}':
		return '{'
	case '<', '>':
		return '<'
	}
	return 0
}

type lineStatus struct {
	corrupted    bool
	errorScore   int64
	incomplete   bool
	stackToClose []bracket
}

func analyzeLine(line string) (lineStatus, error) {
	var stack []bracket

	for _, ch := range line {
		if ch == ' ' || ch == '\n' || ch == '\r' {
			continue
		}

		if isOpen(ch) {
			stack = append(stack, bracket(ch))
		} else if isClose(ch) {
			if len(stack) == 0 {
				return lineStatus{corrupted: true, errorScore: bracket(ch).errorScore()}, nil
			}

			open := getOpenBracket(ch)
			if stack[len(stack)-1] != open {
				return lineStatus{corrupted: true, errorScore: bracket(ch).errorScore()}, nil
			}

			stack = stack[:len(stack)-1]
		} else {
			return lineStatus{}, fmt.Errorf("invalid character: %c", ch)
		}
	}

	if len(stack) > 0 {
		return lineStatus{incomplete: true, stackToClose: stack}, nil
	}

	return lineStatus{}, fmt.Errorf("complete line should not appear")
}

func solve(input string) (int64, int64, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	var part1 int64
	var completionScores []int64

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		status, err := analyzeLine(line)
		if err != nil {
			return 0, 0, fmt.Errorf("error analyzing line %q: %w", line, err)
		}

		if status.corrupted {
			part1 += status.errorScore
		} else if status.incomplete {
			score := int64(0)
			for i := len(status.stackToClose) - 1; i >= 0; i-- {
				br := status.stackToClose[i]
				score = score*5 + br.completionScore()
			}
			completionScores = append(completionScores, score)
		}
	}

	sort.Slice(completionScores, func(i, j int) bool {
		return completionScores[i] < completionScores[j]
	})

	if len(completionScores) == 0 {
		return 0, 0, fmt.Errorf("no incomplete lines found")
	}

	part2 := completionScores[len(completionScores)/2]

	return part1, part2, nil
}
