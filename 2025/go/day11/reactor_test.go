package main

import (
	"errors"
	"os"
	"testing"
)

const officialPart2Example = `svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out`

func TestOfficialPart1Example(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, _, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 5 {
		t.Errorf("part 1 = %d, want 5", part1)
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"official example visits fft before dac", officialPart2Example, 2},
		{"dac before fft", "svr: a b\na: dac\nb: dac\ndac: fft\nfft: out", 2},
		{"only one waypoint on the way", "svr: dac\ndac: out\nfft: out", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part2 != tc.want {
				t.Errorf("part 2 = %d, want %d", part2, tc.want)
			}
		})
	}
}

func TestSinglePath(t *testing.T) {
	part1, _, err := solve("you: out")
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 1 {
		t.Errorf("part 1 = %d, want 1", part1)
	}
}

func TestMalformedInput(t *testing.T) {
	if _, _, err := solve("you out"); err == nil {
		t.Error("want an error for a line without a colon")
	}
	if _, _, err := solve(" \n"); !errors.Is(err, errEmptyInput) {
		t.Errorf("empty input: got %v, want errEmptyInput", err)
	}
}
