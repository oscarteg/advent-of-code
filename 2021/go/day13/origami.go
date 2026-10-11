package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidOrigami = errors.New("invalid origami input")

type point struct {
	x, y int
}

type fold struct {
	axis  byte
	value int
}

func parsePoint(line string) (point, error) {
	x, y, ok := strings.Cut(line, ",")
	if !ok {
		return point{}, fmt.Errorf("%w: malformed dot %q", errInvalidOrigami, line)
	}
	px, err := strconv.Atoi(strings.TrimSpace(x))
	if err != nil {
		return point{}, fmt.Errorf("dot coordinate: %w", err)
	}
	py, err := strconv.Atoi(strings.TrimSpace(y))
	if err != nil {
		return point{}, fmt.Errorf("dot coordinate: %w", err)
	}
	return point{px, py}, nil
}

func parseFold(line string) (fold, error) {
	instruction, ok := strings.CutPrefix(line, "fold along ")
	if !ok {
		return fold{}, fmt.Errorf("%w: malformed fold %q", errInvalidOrigami, line)
	}
	axis, value, ok := strings.Cut(instruction, "=")
	if !ok || (axis != "x" && axis != "y") {
		return fold{}, fmt.Errorf("%w: malformed fold %q", errInvalidOrigami, line)
	}
	v, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fold{}, fmt.Errorf("fold value: %w", err)
	}
	return fold{axis[0], v}, nil
}

func parse(input string) (map[point]struct{}, []fold, error) {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	sections := strings.SplitN(strings.TrimSpace(normalized), "\n\n", 2)
	if len(sections) != 2 {
		return nil, nil, fmt.Errorf("%w: missing fold instructions", errInvalidOrigami)
	}

	dots := make(map[point]struct{})
	for _, line := range strings.Split(sections[0], "\n") {
		p, err := parsePoint(line)
		if err != nil {
			return nil, nil, err
		}
		dots[p] = struct{}{}
	}
	if len(dots) == 0 {
		return nil, nil, fmt.Errorf("%w: no dots given", errInvalidOrigami)
	}

	var folds []fold
	for _, line := range strings.Split(sections[1], "\n") {
		f, err := parseFold(line)
		if err != nil {
			return nil, nil, err
		}
		folds = append(folds, f)
	}
	if len(folds) == 0 {
		return nil, nil, fmt.Errorf("%w: no fold instructions given", errInvalidOrigami)
	}

	return dots, folds, nil
}

func applyFold(dots map[point]struct{}, f fold) map[point]struct{} {
	folded := make(map[point]struct{}, len(dots))
	for p := range dots {
		switch {
		case f.axis == 'x' && p.x > f.value:
			p.x = 2*f.value - p.x
		case f.axis == 'y' && p.y > f.value:
			p.y = 2*f.value - p.y
		}
		folded[p] = struct{}{}
	}
	return folded
}

func render(dots map[point]struct{}) string {
	width, height := 0, 0
	for p := range dots {
		if p.x+1 > width {
			width = p.x + 1
		}
		if p.y+1 > height {
			height = p.y + 1
		}
	}
	rows := make([]string, height)
	for y := 0; y < height; y++ {
		row := make([]byte, width)
		for x := 0; x < width; x++ {
			if _, ok := dots[point{x, y}]; ok {
				row[x] = '#'
			} else {
				row[x] = '.'
			}
		}
		rows[y] = string(row)
	}
	return strings.Join(rows, "\n")
}

func solve(input string) (int, string, error) {
	dots, folds, err := parse(input)
	if err != nil {
		return 0, "", err
	}
	dots = applyFold(dots, folds[0])
	part1 := len(dots)
	for _, f := range folds[1:] {
		dots = applyFold(dots, f)
	}
	return part1, render(dots), nil
}
