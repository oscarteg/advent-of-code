package main

import (
	"testing"
)

const exampleInput = `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2`

func TestSolveExample(t *testing.T) {
	part1, part2, err := solve(exampleInput)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 5 {
		t.Errorf("part1: expected 5, got %d", part1)
	}

	if part2 != 12 {
		t.Errorf("part2: expected 12, got %d", part2)
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		input string
		want  line
		err   bool
	}{
		{
			"0,9 -> 5,9",
			line{point{0, 9}, point{5, 9}},
			false,
		},
		{
			"8,0 -> 0,8",
			line{point{8, 0}, point{0, 8}},
			false,
		},
		{
			"1,1 -> 1,3",
			line{point{1, 1}, point{1, 3}},
			false,
		},
		{
			"invalid",
			line{},
			true,
		},
		{
			"1,2,3 -> 4,5,6",
			line{},
			true,
		},
	}

	for _, tt := range tests {
		l, err := parseLine(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("parseLine(%q) error = %v, want error %v", tt.input, err, tt.err)
			continue
		}
		if err == nil && l != tt.want {
			t.Errorf("parseLine(%q) = %+v, want %+v", tt.input, l, tt.want)
		}
	}
}

func TestLinePoints(t *testing.T) {
	tests := []struct {
		line  line
		count int
		want  []point
	}{
		{
			line{point{0, 9}, point{5, 9}},
			6,
			[]point{{0, 9}, {1, 9}, {2, 9}, {3, 9}, {4, 9}, {5, 9}},
		},
		{
			line{point{1, 1}, point{1, 3}},
			3,
			[]point{{1, 1}, {1, 2}, {1, 3}},
		},
		{
			line{point{9, 7}, point{7, 7}},
			3,
			[]point{{9, 7}, {8, 7}, {7, 7}},
		},
		{
			line{point{0, 0}, point{2, 2}},
			3,
			[]point{{0, 0}, {1, 1}, {2, 2}},
		},
	}

	for _, tt := range tests {
		points := tt.line.points()
		if len(points) != tt.count {
			t.Errorf("line.points() length = %d, want %d", len(points), tt.count)
		}
		for i, p := range points {
			if i < len(tt.want) && p != tt.want[i] {
				t.Errorf("line.points()[%d] = %+v, want %+v", i, p, tt.want[i])
			}
		}
	}
}

func TestIsHorizontalOrVertical(t *testing.T) {
	tests := []struct {
		line line
		want bool
	}{
		{line{point{0, 9}, point{5, 9}}, true},
		{line{point{1, 1}, point{1, 3}}, true},
		{line{point{0, 0}, point{2, 2}}, false},
		{line{point{9, 7}, point{7, 9}}, false},
	}

	for _, tt := range tests {
		if got := tt.line.isHorizontalOrVertical(); got != tt.want {
			t.Errorf("line%+v.isHorizontalOrVertical() = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Errorf("solve empty input: expected error, got nil")
	}
}

func TestMalformedInput(t *testing.T) {
	malformed := "not a valid line"
	_, _, err := solve(malformed)
	if err == nil {
		t.Errorf("solve malformed input: expected error, got nil")
	}
}
