package main

import (
	"os"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 152 || part2 != 301 {
		t.Fatalf("got (%d, %d, %v), want (152, 301, nil)", part1, part2, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name, input  string
		part1, part2 int64
		wantError    bool
	}{
		{
			name:  "human is a direct operand of root",
			input: "root: humn + other\nother: 7\nhumn: 3\n",
			part1: 10,
			part2: 7,
		},
		{
			name:  "inverts operations with the human on the right",
			input: "root: target + chain\ntarget: 10\nchain: base - inner\nbase: 30\ninner: fixed / humn\nfixed: 40\nhumn: 1\n",
			part1: 0,
			part2: 2,
		},
		{
			name:  "inverts multiplication and division with the human on the left",
			input: "root: a + b\na: 36\nb: c * three\nthree: 3\nc: humn / two\ntwo: 2\nhumn: 0\n",
			part1: 36,
			part2: 24,
		},
		{
			name:  "ignores blank lines and surrounding whitespace",
			input: "\n  root: humn + other  \r\n\nother: 7\n humn: 3\n\n",
			part1: 10,
			part2: 7,
		},
		{name: "empty input", input: "", wantError: true},
		{name: "missing colon", input: "root 1\n", wantError: true},
		{name: "missing name", input: ": 1\n", wantError: true},
		{name: "missing job", input: "root:\n", wantError: true},
		{name: "too many operands", input: "root: a + b + c\n", wantError: true},
		{name: "unknown operator", input: "root: a ? b\n", wantError: true},
		{name: "number is not a number", input: "root: humn + other\nother: twelve\nhumn: 1\n", wantError: true},
		{name: "duplicate monkey", input: "root: humn + other\nother: 7\nother: 8\nhumn: 1\n", wantError: true},
		{name: "unknown monkey", input: "root: humn + missing\nhumn: 1\n", wantError: true},
		{name: "missing human", input: "root: a + b\na: 1\nb: 2\n", wantError: true},
		{name: "cycle", input: "root: a + b\na: 1\nb: b - a\nhumn: 1\n", wantError: true},
		{name: "no integer answer", input: "root: odd + scaled\nodd: 7\nscaled: humn * two\ntwo: 2\nhumn: 0\n", wantError: true},
		{name: "division by zero", input: "root: a + b\na: 1\nb: a / zero\nzero: 0\nhumn: 1\n", wantError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
