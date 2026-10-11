package main

import (
	"fmt"
	"strconv"
	"strings"
)

type card struct {
	winningNumbers map[int]struct{}
	yourNumbers    []int
}

func parseCard(line string) (card, error) {
	_, numbers, ok := strings.Cut(line, ":")
	if !ok {
		return card{}, fmt.Errorf("card missing ':': %q", line)
	}
	winningField, yourField, ok := strings.Cut(numbers, "|")
	if !ok {
		return card{}, fmt.Errorf("card missing '|': %q", line)
	}
	winningNumbers := make(map[int]struct{})
	for _, field := range strings.Fields(winningField) {
		number, err := strconv.Atoi(field)
		if err != nil {
			return card{}, fmt.Errorf("winning number: %w", err)
		}
		winningNumbers[number] = struct{}{}
	}
	yourNumbers := make([]int, 0)
	for _, field := range strings.Fields(yourField) {
		number, err := strconv.Atoi(field)
		if err != nil {
			return card{}, fmt.Errorf("your number: %w", err)
		}
		yourNumbers = append(yourNumbers, number)
	}
	return card{winningNumbers: winningNumbers, yourNumbers: yourNumbers}, nil
}

func parseCards(input string) ([]card, error) {
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n")
	cards := make([]card, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		c, err := parseCard(line)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("no cards found")
	}
	return cards, nil
}

func (c card) matches() int {
	count := 0
	for _, number := range c.yourNumbers {
		if _, ok := c.winningNumbers[number]; ok {
			count++
		}
	}
	return count
}

func solve(input string) (part1, part2 int, err error) {
	cards, err := parseCards(input)
	if err != nil {
		return 0, 0, err
	}

	points := 0
	copies := make([]int, len(cards))
	for i := range copies {
		copies[i] = 1
	}
	for index, c := range cards {
		matches := c.matches()
		if matches > 0 {
			points += 1 << (matches - 1)
		}
		for offset := 1; offset <= matches; offset++ {
			target := index + offset
			if target < len(copies) {
				copies[target] += copies[index]
			}
		}
	}

	total := 0
	for _, count := range copies {
		total += count
	}
	return points, total, nil
}
