package day13

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	s, err := parseSchedule(input)
	if err != nil {
		return 0, err
	}

	bestID, bestWait := 0, -1
	for _, b := range s.buses {
		wait := (b.id - s.earliest%b.id) % b.id
		if bestWait == -1 || wait < bestWait {
			bestID, bestWait = b.id, wait
		}
	}

	return bestID * bestWait, nil
}
