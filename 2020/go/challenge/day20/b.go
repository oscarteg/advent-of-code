package day20

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	tiles, err := parseTiles(input)
	if err != nil {
		return 0, err
	}

	owners := ownersByEdge(tiles)
	corners, err := cornerIndices(tiles, owners)
	if err != nil {
		return 0, err
	}

	assembled, err := assemble(tiles, owners, corners)
	if err != nil {
		return 0, err
	}
	return waterRoughness(assembled), nil
}
