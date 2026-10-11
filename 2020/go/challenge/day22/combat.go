package day22

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func parseDecks(input *challenge.Input) ([]int, []int, error) {
	lines := input.LinesSlice()

	var decks [][]int
	var current []int
	started := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if started {
				decks = append(decks, current)
				current = nil
				started = false
			}
			continue
		}

		if strings.HasPrefix(line, "Player") && strings.HasSuffix(line, ":") {
			started = true
			continue
		}

		card, err := strconv.Atoi(line)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid card %q: %w", line, err)
		}
		current = append(current, card)
	}
	if started {
		decks = append(decks, current)
	}

	if len(decks) != 2 {
		return nil, nil, fmt.Errorf("expected exactly 2 player decks, got %d", len(decks))
	}
	if len(decks[0]) == 0 || len(decks[1]) == 0 {
		return nil, nil, fmt.Errorf("both players need at least one card")
	}

	return decks[0], decks[1], nil
}

func score(deck []int) int {
	total := 0
	n := len(deck)
	for i, card := range deck {
		total += card * (n - i)
	}
	return total
}

func playCombat(player1, player2 []int) int {
	p1, p2 := append([]int(nil), player1...), append([]int(nil), player2...)

	for len(p1) > 0 && len(p2) > 0 {
		card1, card2 := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]

		if card1 > card2 {
			p1 = append(p1, card1, card2)
		} else {
			p2 = append(p2, card2, card1)
		}
	}

	if len(p1) == 0 {
		return score(p2)
	}
	return score(p1)
}

func deckKey(player1, player2 []int) string {
	var b strings.Builder
	for _, c := range player1 {
		fmt.Fprintf(&b, "%d,", c)
	}
	b.WriteByte('|')
	for _, c := range player2 {
		fmt.Fprintf(&b, "%d,", c)
	}
	return b.String()
}

// playRecursiveCombat plays Recursive Combat and reports whether player 1 won
// along with the winner's final deck.
func playRecursiveCombat(player1, player2 []int) (bool, []int) {
	p1, p2 := append([]int(nil), player1...), append([]int(nil), player2...)
	seen := make(map[string]bool)

	for len(p1) > 0 && len(p2) > 0 {
		key := deckKey(p1, p2)
		if seen[key] {
			return true, p1
		}
		seen[key] = true

		card1, card2 := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]

		var player1WinsRound bool
		if card1 <= len(p1) && card2 <= len(p2) {
			player1WinsRound, _ = playRecursiveCombat(p1[:card1], p2[:card2])
		} else {
			player1WinsRound = card1 > card2
		}

		if player1WinsRound {
			p1 = append(p1, card1, card2)
		} else {
			p2 = append(p2, card2, card1)
		}
	}

	if len(p1) == 0 {
		return false, p2
	}
	return true, p1
}
