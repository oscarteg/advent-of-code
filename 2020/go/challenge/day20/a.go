package day20

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	tiles, err := parseTiles(input)
	if err != nil {
		return 0, err
	}

	corners, err := cornerIndices(tiles, ownersByEdge(tiles))
	if err != nil {
		return 0, err
	}

	product := 1
	for _, index := range corners {
		product *= tiles[index].id
	}
	return product, nil
}
