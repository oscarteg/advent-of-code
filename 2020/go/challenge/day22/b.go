package day22

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	player1, player2, err := parseDecks(input)
	if err != nil {
		return 0, err
	}

	_, winningDeck := playRecursiveCombat(player1, player2)
	return score(winningDeck), nil
}
