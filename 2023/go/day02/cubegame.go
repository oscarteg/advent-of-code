package main

import (
	"fmt"
	"strconv"
	"strings"
)

type reveal struct {
	red, green, blue uint64
}

type game struct {
	id      uint64
	reveals []reveal
}

func (g game) isPossible(red, green, blue uint64) bool {
	for _, r := range g.reveals {
		if r.red > red || r.green > green || r.blue > blue {
			return false
		}
	}
	return true
}

func (g game) power() uint64 {
	var max reveal
	for _, r := range g.reveals {
		max.red = maxUint64(max.red, r.red)
		max.green = maxUint64(max.green, r.green)
		max.blue = maxUint64(max.blue, r.blue)
	}
	return max.red * max.green * max.blue
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func parseGame(line string) (game, error) {
	header, body, found := strings.Cut(line, ":")
	if !found {
		return game{}, fmt.Errorf("parse game: missing ':' in %q", line)
	}
	idField, found := strings.CutPrefix(strings.TrimSpace(header), "Game ")
	if !found {
		return game{}, fmt.Errorf("parse game: missing 'Game ' prefix in %q", header)
	}
	id, err := strconv.ParseUint(idField, 10, 64)
	if err != nil {
		return game{}, fmt.Errorf("parse game id: %w", err)
	}

	subsets := strings.Split(body, ";")
	reveals := make([]reveal, 0, len(subsets))
	for _, subset := range subsets {
		r, err := parseReveal(subset)
		if err != nil {
			return game{}, err
		}
		reveals = append(reveals, r)
	}
	if len(reveals) == 0 {
		return game{}, fmt.Errorf("parse game: no reveals in %q", line)
	}
	return game{id: id, reveals: reveals}, nil
}

func parseReveal(subset string) (reveal, error) {
	var r reveal
	for _, cube := range strings.Split(subset, ",") {
		count, color, found := strings.Cut(strings.TrimSpace(cube), " ")
		if !found {
			return reveal{}, fmt.Errorf("parse cube: missing count and color in %q", cube)
		}
		n, err := strconv.ParseUint(count, 10, 64)
		if err != nil {
			return reveal{}, fmt.Errorf("parse cube count: %w", err)
		}
		switch color {
		case "red":
			r.red = n
		case "green":
			r.green = n
		case "blue":
			r.blue = n
		default:
			return reveal{}, fmt.Errorf("parse cube: unknown color %q", color)
		}
	}
	return r, nil
}

func parseGames(input string) ([]game, error) {
	var games []game
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		g, err := parseGame(line)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	if len(games) == 0 {
		return nil, fmt.Errorf("parse games: no games found")
	}
	return games, nil
}

const (
	bagRed   = 12
	bagGreen = 13
	bagBlue  = 14
)

func solve(input string) (uint64, uint64, error) {
	games, err := parseGames(input)
	if err != nil {
		return 0, 0, err
	}
	return solveWith(games, bagRed, bagGreen, bagBlue)
}

func solveWith(games []game, red, green, blue uint64) (uint64, uint64, error) {
	var possibleSum, powerSum uint64
	for _, g := range games {
		if g.isPossible(red, green, blue) {
			possibleSum += g.id
		}
		powerSum += g.power()
	}
	return possibleSum, powerSum, nil
}
