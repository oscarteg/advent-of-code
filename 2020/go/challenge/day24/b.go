package day24

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	tiles := make(map[pos]int)

	for line := range input.Lines() {
		if line == "" {
			continue
		}

		p, err := parsePath(line)
		if err != nil {
			return 0, err
		}

		tiles[p]++
	}

	blackTiles := make(map[pos]bool)
	for p, flips := range tiles {
		if flips%2 == 1 {
			blackTiles[p] = true
		}
	}

	for day := 0; day < 100; day++ {
		blackTiles = simulateOneDay(blackTiles)
	}

	return len(blackTiles), nil
}

func countBlackNeighbors(p pos, blackTiles map[pos]bool) int {
	count := 0
	for _, d := range []pos{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, -1}, {-1, 1},
	} {
		neighbor := pos{p.x + d.x, p.y + d.y}
		if blackTiles[neighbor] {
			count++
		}
	}
	return count
}

func simulateOneDay(blackTiles map[pos]bool) map[pos]bool {
	toCheck := make(map[pos]bool)

	for p := range blackTiles {
		toCheck[p] = true
		for _, d := range []pos{
			{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, -1}, {-1, 1},
		} {
			neighbor := pos{p.x + d.x, p.y + d.y}
			toCheck[neighbor] = true
		}
	}

	newBlackTiles := make(map[pos]bool)
	for p := range toCheck {
		isBlack := blackTiles[p]
		blackNeighbors := countBlackNeighbors(p, blackTiles)

		shouldBeBlack := false
		if isBlack {
			shouldBeBlack = blackNeighbors == 1 || blackNeighbors == 2
		} else {
			shouldBeBlack = blackNeighbors == 2
		}

		if shouldBeBlack {
			newBlackTiles[p] = true
		}
	}

	return newBlackTiles
}
