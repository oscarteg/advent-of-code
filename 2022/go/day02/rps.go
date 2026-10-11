package main

import (
	"fmt"
	"strings"
)

type Move int

const (
	Rock     Move = 1
	Paper    Move = 2
	Scissors Move = 3
)

type Outcome int

const (
	Loss Outcome = 0
	Draw Outcome = 3
	Win  Outcome = 6
)

func parseMove(ch rune) (Move, error) {
	switch ch {
	case 'A', 'X':
		return Rock, nil
	case 'B', 'Y':
		return Paper, nil
	case 'C', 'Z':
		return Scissors, nil
	default:
		return 0, fmt.Errorf("invalid move: %c", ch)
	}
}

func parseOutcome(ch rune) (Outcome, error) {
	switch ch {
	case 'X':
		return Loss, nil
	case 'Y':
		return Draw, nil
	case 'Z':
		return Win, nil
	default:
		return 0, fmt.Errorf("invalid outcome: %c", ch)
	}
}

func moveThatDefeats(opponent Move) Move {
	switch opponent {
	case Rock:
		return Paper
	case Paper:
		return Scissors
	case Scissors:
		return Rock
	default:
		return Rock
	}
}

func moveThatLosesTo(opponent Move) Move {
	switch opponent {
	case Rock:
		return Scissors
	case Paper:
		return Rock
	case Scissors:
		return Paper
	default:
		return Rock
	}
}

func scoreRound(ourMove Move, opponent Move) int {
	var outcome Outcome
	if ourMove == opponent {
		outcome = Draw
	} else if (ourMove == Rock && opponent == Scissors) ||
		(ourMove == Paper && opponent == Rock) ||
		(ourMove == Scissors && opponent == Paper) {
		outcome = Win
	} else {
		outcome = Loss
	}
	return int(ourMove) + int(outcome)
}

func solve(input string) (int, int, error) {
	var part1, part2 int

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		if len(line) < 3 {
			return 0, 0, fmt.Errorf("invalid line format: %q", line)
		}

		opponentChar := rune(line[0])
		codeChar := rune(line[2])

		opponent, err := parseMove(opponentChar)
		if err != nil {
			return 0, 0, err
		}

		ourMovePart1, err := parseMove(codeChar)
		if err != nil {
			return 0, 0, err
		}
		part1 += scoreRound(ourMovePart1, opponent)

		outcome, err := parseOutcome(codeChar)
		if err != nil {
			return 0, 0, err
		}
		var ourMovePart2 Move
		switch outcome {
		case Loss:
			ourMovePart2 = moveThatLosesTo(opponent)
		case Draw:
			ourMovePart2 = opponent
		case Win:
			ourMovePart2 = moveThatDefeats(opponent)
		}
		part2 += scoreRound(ourMovePart2, opponent)
	}

	return part1, part2, nil
}
