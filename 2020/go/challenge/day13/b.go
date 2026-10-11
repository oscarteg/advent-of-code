package day13

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	s, err := parseSchedule(input)
	if err != nil {
		return 0, err
	}

	timestamp, period := 0, 1
	for _, bus := range s.buses {
		for (timestamp+bus.offset)%bus.id != 0 {
			timestamp += period
		}
		period *= bus.id
	}

	return timestamp, nil
}
