package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var errInvalidBingo = errors.New("invalid bingo game")

type board struct {
	numbers [25]uint64
	marked  [25]bool
	won     bool
}

func parseBoard(input string) (board, error) {
	var result board
	rows := strings.Split(strings.TrimSpace(input), "\n")
	if len(rows) != 5 {
		return result, fmt.Errorf("%w: board needs five rows", errInvalidBingo)
	}
	for row, line := range rows {
		fields := strings.Fields(line)
		if len(fields) != 5 {
			return result, fmt.Errorf("%w: row needs five numbers", errInvalidBingo)
		}
		for column, field := range fields {
			number, err := strconv.ParseUint(field, 10, 32)
			if err != nil {
				return result, fmt.Errorf("board number: %w", err)
			}
			result.numbers[row*5+column] = number
		}
	}
	return result, nil
}

func (b *board) mark(draw uint64) bool {
	for i, number := range b.numbers {
		b.marked[i] = b.marked[i] || number == draw
	}
	for line := 0; line < 5; line++ {
		row, column := true, true
		for offset := 0; offset < 5; offset++ {
			row = row && b.marked[line*5+offset]
			column = column && b.marked[offset*5+line]
		}
		if row || column {
			b.won = true
			return true
		}
	}
	return false
}

func solve(input string) (uint64, uint64, error) {
	sections := strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n\n")
	if len(sections) < 2 || strings.TrimSpace(sections[0]) == "" {
		return 0, 0, fmt.Errorf("%w: missing boards", errInvalidBingo)
	}
	draws := make([]uint64, 0)
	for _, field := range strings.Split(sections[0], ",") {
		number, err := strconv.ParseUint(strings.TrimSpace(field), 10, 32)
		if err != nil {
			return 0, 0, fmt.Errorf("draw number: %w", err)
		}
		draws = append(draws, number)
	}
	if len(draws) == 0 {
		return 0, 0, fmt.Errorf("%w: missing draws", errInvalidBingo)
	}
	boards := make([]board, 0, len(sections)-1)
	for _, section := range sections[1:] {
		b, err := parseBoard(section)
		if err != nil {
			return 0, 0, err
		}
		boards = append(boards, b)
	}
	remaining := len(boards)
	var first uint64
	for _, draw := range draws {
		for i := range boards {
			b := &boards[i]
			if b.won || !b.mark(draw) {
				continue
			}
			var unmarked uint64
			for index, number := range b.numbers {
				if !b.marked[index] {
					unmarked += number
				}
			}
			if draw > 0 && unmarked > math.MaxUint64/draw {
				return 0, 0, fmt.Errorf("%w: score overflows", errInvalidBingo)
			}
			score := unmarked * draw
			if remaining == len(boards) {
				first = score
			}
			remaining--
			if remaining == 0 {
				return first, score, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("%w: draws end before every board wins", errInvalidBingo)
}
