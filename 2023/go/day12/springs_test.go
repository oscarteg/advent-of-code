package main

import (
	"errors"
	"os"
	"strconv"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 21 || part2 != 525152 {
		t.Fatalf("got (%d, %d, %v), want (21, 525152, nil)", part1, part2, err)
	}
}

func TestSingleCopyLeavesTheRowsFolded(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solveWith(string(input), 1)

	if err != nil || part1 != 21 || part2 != 21 {
		t.Fatalf("got (%d, %d, %v), want (21, 21, nil)", part1, part2, err)
	}
}

func TestRowArrangements(t *testing.T) {
	tests := []struct {
		name     string
		row      string
		folded   int
		unfolded int
	}{
		{"single arrangement", "???.### 1,1,3", 1, 1},
		{"two pairs of unknowns", ".??..??...?##. 1,1,3", 4, 16384},
		{"alternating unknowns", "?#?#?#?#?#?#?#? 1,3,1,6", 1, 1},
		{"known groups behind unknowns", "????.#...#... 4,1,1", 1, 16},
		{"long known groups", "????.######..#####. 1,6,5", 4, 2500},
		{"mostly unknown", "?###???????? 3,2,1", 10, 506250},
		{"no unknown springs", "##.# 2,1", 1, 1},
		{"known row that cannot fit", "#.# 2", 0, 0},
		{"groups fill the whole row", "????? 5", 1, 1},
		{"wide open row", "????? 1", 5, 53130},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.row)
			if err != nil {
				t.Fatal(err)
			}
			if part1 != tc.folded || part2 != tc.unfolded {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.folded, tc.unfolded)
			}
		})
	}
}

func TestUnfoldJoinsCopiesWithUnknown(t *testing.T) {
	unfolded := record{springs: ".#", groups: []int{1}}.unfold(5)

	if unfolded.springs != ".#?.#?.#?.#?.#" {
		t.Fatalf("got springs %q, want %q", unfolded.springs, ".#?.#?.#?.#?.#")
	}
	if len(unfolded.groups) != 5 {
		t.Fatalf("got groups %v, want five sizes", unfolded.groups)
	}
	for _, group := range unfolded.groups {
		if group != 1 {
			t.Fatalf("got groups %v, want every size to be 1", unfolded.groups)
		}
	}
}

func TestMalformedInputReturnsErrors(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty input", ""},
		{"missing group sizes", "???.###"},
		{"extra field", "???.### 1,1,3 4"},
		{"unknown symbol", "??!.### 1,1,3"},
		{"group size is not a number", "???.### 1,x,3"},
		{"zero group size", "???.### 1,0,3"},
		{"negative group size", "???.### -1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("expected an error for %q", tc.input)
			}
		})
	}
}

func TestMalformedInputWrapsCauses(t *testing.T) {
	if _, _, err := solve("??x 1"); !errors.Is(err, errInvalidRecords) {
		t.Fatalf("got %v, want errInvalidRecords", err)
	}
	if _, _, err := solve("??? 1,x"); !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("got %v, want strconv.ErrSyntax", err)
	}
}
