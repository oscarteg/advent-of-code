package main

import (
	"os"
	"testing"
)

func TestSolveOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("reading example: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 95437 {
		t.Errorf("part1 = %d, want 95437", part1)
	}
	if part2 != 24933642 {
		t.Errorf("part2 = %d, want 24933642", part2)
	}
}

func TestSolveEmptyDirectoriesCountAsZero(t *testing.T) {
	input := "$ cd /\n$ ls\ndir a\n100000 big.txt\n$ cd a\n$ ls\n"
	part1, _, err := solve(input)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 100000 {
		t.Errorf("part1 = %d, want 100000", part1)
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	cases := map[string]string{
		"cd above root":        "$ cd /\n$ cd ..\n",
		"malformed listing":    "$ cd /\n$ ls\nnotanumber\n",
		"empty directory name": "$ cd /\n$ ls\ndir \n",
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, _, err := solve(input); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}
