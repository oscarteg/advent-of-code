package main

import (
	"os"
	"strings"
	"testing"
)

const smallExample = `on x=10..12,y=10..12,z=10..12
on x=11..13,y=11..13,z=11..13
off x=9..11,y=9..11,z=9..11
on x=10..10,y=10..10,z=10..10
`

const largerExample = `on x=-20..26,y=-36..17,z=-47..7
on x=-20..33,y=-21..23,z=-26..28
on x=-22..28,y=-29..23,z=-38..16
on x=-46..7,y=-6..46,z=-50..-1
on x=-49..1,y=-3..46,z=-24..28
on x=2..47,y=-22..22,z=-23..27
on x=-27..23,y=-28..26,z=-21..29
on x=-39..5,y=-6..47,z=-3..44
on x=-30..21,y=-8..43,z=-13..34
on x=-22..26,y=-27..20,z=-29..19
off x=-48..-32,y=26..41,z=-47..-37
on x=-12..35,y=6..50,z=-50..-2
off x=-48..-32,y=-32..-16,z=-15..-5
on x=-18..26,y=-33..15,z=-7..46
off x=-40..-22,y=-38..-28,z=23..41
on x=-16..35,y=-41..10,z=-47..6
off x=-32..-23,y=11..30,z=-14..3
on x=-49..-5,y=-3..45,z=-29..18
off x=18..30,y=-20..-8,z=-3..13
on x=-41..9,y=-7..43,z=-33..15
on x=-54112..-39298,y=-85059..-49293,z=-27449..7877
on x=967..23432,y=45373..81175,z=27513..53682
`

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 474140 || part2 != 2758514936282235 {
		t.Fatalf("got (%d, %d, %v), want (474140, 2758514936282235, nil)", part1, part2, err)
	}
}

func TestLargerExample(t *testing.T) {
	part1, _, err := solve(largerExample)

	if err != nil || part1 != 590784 {
		t.Fatalf("got (%d, %v), want (590784, nil)", part1, err)
	}
}

func TestRebootSteps(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
	}{
		{"small example", smallExample, 39, 39},
		{"windows line endings", strings.ReplaceAll(smallExample, "\n", "\r\n"), 39, 39},
		{"outside the initialization region", "on x=100..101,y=100..101,z=100..101", 0, 8},
		{"clipped at the region border", "on x=49..51,y=0..0,z=0..0", 2, 3},
		{"overlap counted once", "on x=0..1,y=0..0,z=0..0\non x=1..2,y=0..0,z=0..0", 3, 3},
		{"off clears earlier cubes", "on x=0..1,y=0..1,z=0..1\noff x=0..1,y=0..1,z=0..1", 0, 0},
		{"off without cubes on", "off x=0..1,y=0..1,z=0..1", 0, 0},
		{"empty input", "", 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part1 != tc.part1 || part2 != tc.part2 {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"missing state", "x=0..1,y=0..1,z=0..1"},
		{"unknown state", "maybe x=0..1,y=0..1,z=0..1"},
		{"missing axis", "on x=0..1,y=0..1"},
		{"wrong axis order", "on y=0..1,x=0..1,z=0..1"},
		{"range without dots", "on x=0,y=0..1,z=0..1"},
		{"non numeric coordinate", "on x=a..1,y=0..1,z=0..1"},
		{"reversed range", "on x=1..0,y=0..1,z=0..1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want error", tc.input)
			}
		})
	}
}
