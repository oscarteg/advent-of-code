package day22

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	player1, player2, err := parseDecks(input)
	if err != nil {
		return 0, err
	}

	return playCombat(player1, player2), nil
}
