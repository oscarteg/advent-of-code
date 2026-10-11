package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidGame = errors.New("invalid dirac dice game")

type player struct {
	position uint8
	score    uint16
}

func parsePlayer(line string, expected int) (player, error) {
	label, value, ok := strings.Cut(line, ": ")
	if !ok {
		return player{}, fmt.Errorf("%w: malformed starting position line %q", errInvalidGame, line)
	}
	if strings.TrimSpace(label) != fmt.Sprintf("Player %d starting position", expected) {
		return player{}, fmt.Errorf("%w: expected player %d, got %q", errInvalidGame, expected, label)
	}
	position, err := strconv.ParseUint(strings.TrimSpace(value), 10, 8)
	if err != nil {
		return player{}, fmt.Errorf("starting position: %w", err)
	}
	if position < 1 || position > 10 {
		return player{}, fmt.Errorf("%w: starting position must be between 1 and 10", errInvalidGame)
	}
	return player{position: uint8(position)}, nil
}

func parse(input string) (player, player, error) {
	lines := make([]string, 0, 2)
	for _, line := range strings.Split(input, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) != 2 {
		return player{}, player{}, fmt.Errorf("%w: expected two starting position lines", errInvalidGame)
	}
	player1, err := parsePlayer(lines[0], 1)
	if err != nil {
		return player{}, player{}, err
	}
	player2, err := parsePlayer(lines[1], 2)
	if err != nil {
		return player{}, player{}, err
	}
	return player1, player2, nil
}

func advance(position uint8, roll uint32) uint8 {
	return uint8((uint32(position-1)+roll)%10 + 1)
}

func playDeterministic(players [2]player) uint64 {
	die := uint32(0)
	rolls := uint64(0)
	turn := 0
	for {
		total := uint32(0)
		for i := 0; i < 3; i++ {
			die = die%100 + 1
			total += die
		}
		rolls += 3
		players[turn].position = advance(players[turn].position, total)
		players[turn].score += uint16(players[turn].position)
		if players[turn].score >= 1000 {
			loser := players[1-turn]
			return uint64(loser.score) * rolls
		}
		turn = 1 - turn
	}
}

type rollSum struct {
	total        uint32
	multiplicity uint64
}

var rollSums = [7]rollSum{
	{3, 1}, {4, 3}, {5, 6}, {6, 7}, {7, 6}, {8, 3}, {9, 1},
}

type gameState struct {
	current, other player
}

func countWins(current, other player, memo map[gameState][2]uint64) (uint64, uint64) {
	if other.score >= 21 {
		return 0, 1
	}
	state := gameState{current, other}
	if wins, ok := memo[state]; ok {
		return wins[0], wins[1]
	}
	var winsCurrent, winsOther uint64
	for _, roll := range rollSums {
		position := advance(current.position, roll.total)
		moved := player{position: position, score: current.score + uint16(position)}
		otherWins, currentWins := countWins(other, moved, memo)
		winsCurrent += roll.multiplicity * currentWins
		winsOther += roll.multiplicity * otherWins
	}
	memo[state] = [2]uint64{winsCurrent, winsOther}
	return winsCurrent, winsOther
}

func playDirac(players [2]player) uint64 {
	memo := make(map[gameState][2]uint64)
	wins1, wins2 := countWins(players[0], players[1], memo)
	if wins1 > wins2 {
		return wins1
	}
	return wins2
}

func solve(input string) (uint64, uint64, error) {
	player1, player2, err := parse(input)
	if err != nil {
		return 0, 0, err
	}
	part1 := playDeterministic([2]player{player1, player2})
	part2 := playDirac([2]player{player1, player2})
	return part1, part2, nil
}
