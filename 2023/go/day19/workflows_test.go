package main

import (
	"os"
	"testing"
)

// all is every rating combination of 1..4000 in four categories; one is a
// single rating in one category with the other three unconstrained.
const (
	all = 4000 * 4000 * 4000 * 4000
	one = 4000 * 4000 * 4000
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 19114 || part2 != 167409079868000 {
		t.Fatalf("got (%d, %d, %v), want (19114, 167409079868000, nil)", part1, part2, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		part1 int64
		part2 int64
	}{
		{"accept everything", "in{A}\n\n{x=1,m=2,a=3,s=4}", 10, all},
		{"reject everything", "in{R}\n\n{x=1,m=2,a=3,s=4}", 0, 0},
		{"chained workflow", "in{s<2:hop,R}\nhop{A}\n\n{x=1,m=1,a=1,s=1}", 4, one},
		{"comparisons are strict", "in{x<2:A,R}\n\n{x=2,m=1,a=1,s=1}\n{x=1,m=1,a=1,s=1}", 4, one},
		{"upper bound is reachable", "in{x>3999:A,R}\n\n{x=4000,m=1,a=1,s=1}", 4003, one},
		{"first matching rule wins", "in{x>5:R,x>1:A,R}\n\n{x=10,m=1,a=1,s=1}\n{x=3,m=1,a=1,s=1}", 6, 4 * one},
		{"both branches accepted", "in{x<2001:A,A}\n\n{x=1,m=1,a=1,s=1}", 4, all},
		{"unsatisfiable condition", "in{x<1:A,R}\n\n{x=1,m=1,a=1,s=1}", 0, 0},
		{"fall-through keeps the rest", "in{x>10:A,x>5:A,R}\n\n{x=7,m=1,a=1,s=1}", 10, 3995 * one},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if part1 != tc.part1 || part2 != tc.part2 {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	for _, input := range []string{
		"",
		"in{A}",
		"{x=1,m=1,a=1,s=1}",
		"px{A}\n\n{x=1,m=1,a=1,s=1}",
		"in A\n\n{x=1,m=1,a=1,s=1}",
		"in{A\n\n{x=1,m=1,a=1,s=1}",
		"in{A}\nin{R}\n\n{x=1,m=1,a=1,s=1}",
		"in{x<1:A}\n\n{x=1,m=1,a=1,s=1}",
		"in{A,x<1:R}\n\n{x=1,m=1,a=1,s=1}",
		"in{q<1:A,A}\n\n{x=1,m=1,a=1,s=1}",
		"in{x=1:A,A}\n\n{x=1,m=1,a=1,s=1}",
		"in{x<:A,A}\n\n{x=1,m=1,a=1,s=1}",
		"in{x<1:,A}\n\n{x=1,m=1,a=1,s=1}",
		"in{x<1:gone,A}\n\n{x=1,m=1,a=1,s=1}",
		"in{loop}\nloop{in}\n\n{x=1,m=1,a=1,s=1}",
		"in{A}\n\n{x=1,m=1,a=1}",
		"in{A}\n\n{x=1,m=1,a=1,s=1,s=2}",
		"in{A}\n\n{x=1,m=1,a=1,s=x}",
		"in{A}\n\nx=1,m=1,a=1,s=1",
	} {
		if _, _, err := solve(input); err == nil {
			t.Fatalf("solve(%q): want error, got nil", input)
		}
	}
}
