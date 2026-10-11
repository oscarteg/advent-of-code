package day16

import (
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

const departurePrefix = "departure"

func B(input *challenge.Input) (int, error) {
	parsed, err := parseNotes(input.LinesSlice())
	if err != nil {
		return 0, err
	}
	order, err := parsed.fieldOrder()
	if err != nil {
		return 0, err
	}

	// Notes without a departure field (such as the examples) leave an empty product of 1.
	product := 1
	for column, name := range order {
		if strings.HasPrefix(name, departurePrefix) {
			product *= parsed.ticket[column]
		}
	}
	return product, nil
}
