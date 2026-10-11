package main

import (
	"fmt"
	"strings"
)

var errInvalidTrenchMap = fmt.Errorf("invalid trench map")

type point struct {
	x, y int
}

type image struct {
	lit        map[point]bool
	minX, maxX int
	minY, maxY int
}

func parseAlgorithm(raw string) ([512]bool, error) {
	var algorithm [512]bool
	cleaned := strings.Map(func(r rune) rune {
		if r == '#' || r == '.' {
			return r
		}
		return -1
	}, raw)
	if len(cleaned) != 512 {
		return algorithm, fmt.Errorf("%w: algorithm must have 512 entries, got %d", errInvalidTrenchMap, len(cleaned))
	}
	for i, ch := range cleaned {
		algorithm[i] = ch == '#'
	}
	return algorithm, nil
}

func parseImage(raw string) (image, error) {
	rows := strings.Split(strings.Trim(raw, "\n"), "\n")
	if len(rows) == 0 || rows[0] == "" {
		return image{}, fmt.Errorf("%w: missing image rows", errInvalidTrenchMap)
	}
	width := len(rows[0])
	lit := make(map[point]bool)
	for y, row := range rows {
		if len(row) != width {
			return image{}, fmt.Errorf("%w: image rows must have equal length", errInvalidTrenchMap)
		}
		for x, ch := range row {
			switch ch {
			case '#':
				lit[point{x, y}] = true
			case '.':
			default:
				return image{}, fmt.Errorf("%w: invalid image pixel %q", errInvalidTrenchMap, ch)
			}
		}
	}
	return image{lit: lit, minX: 0, maxX: width - 1, minY: 0, maxY: len(rows) - 1}, nil
}

func parseInput(input string) ([512]bool, image, error) {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	algorithmPart, imagePart, found := strings.Cut(normalized, "\n\n")
	if !found {
		return [512]bool{}, image{}, fmt.Errorf("%w: missing blank line between algorithm and image", errInvalidTrenchMap)
	}
	algorithm, err := parseAlgorithm(algorithmPart)
	if err != nil {
		return [512]bool{}, image{}, err
	}
	img, err := parseImage(imagePart)
	if err != nil {
		return [512]bool{}, image{}, err
	}
	return algorithm, img, nil
}

func (img image) pixelAt(x, y int, background bool) bool {
	if x < img.minX || x > img.maxX || y < img.minY || y > img.maxY {
		return background
	}
	return img.lit[point{x, y}]
}

func (img image) enhance(algorithm [512]bool, background bool) (image, bool) {
	minX, minY := img.minX-1, img.minY-1
	maxX, maxY := img.maxX+1, img.maxY+1
	lit := make(map[point]bool)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			index := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					index <<= 1
					if img.pixelAt(x+dx, y+dy, background) {
						index |= 1
					}
				}
			}
			if algorithm[index] {
				lit[point{x, y}] = true
			}
		}
	}
	nextBackgroundIndex := 0
	if background {
		nextBackgroundIndex = 511
	}
	nextBackground := algorithm[nextBackgroundIndex]
	return image{lit: lit, minX: minX, maxX: maxX, minY: minY, maxY: maxY}, nextBackground
}

func enhanceNTimes(img image, algorithm [512]bool, steps int) int {
	background := false
	for i := 0; i < steps; i++ {
		img, background = img.enhance(algorithm, background)
	}
	return len(img.lit)
}

func solveWith(input string, part1Steps, part2Steps int) (int, int, error) {
	algorithm, img, err := parseInput(input)
	if err != nil {
		return 0, 0, err
	}
	part1 := enhanceNTimes(img, algorithm, part1Steps)
	part2 := enhanceNTimes(img, algorithm, part2Steps)
	return part1, part2, nil
}

func solve(input string) (int, int, error) {
	return solveWith(input, 2, 50)
}
