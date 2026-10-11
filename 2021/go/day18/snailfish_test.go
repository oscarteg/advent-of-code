package main

import (
	"os"
	"slices"
	"testing"
)

const largerExample = `[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]
[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]
[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]
[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]
[7,[5,[[3,8],[1,4]]]]
[[2,[2,2]],[8,[8,1]]]
[2,9]
[1,[[[9,3],9],[[9,0],[0,7]]]]
[[[5,[7,4]],7],1]
[[[[4,2],2],6],[8,7]]`

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 4140 || part2 != 3993 {
		t.Fatalf("got (%d, %d, %v), want (4140, 3993, nil)", part1, part2, err)
	}
}

func TestMagnitudes(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"[[1,2],[[3,4],5]]", 143},
		{"[[[[0,7],4],[[7,8],[6,0]]],[8,1]]", 1384},
		{"[[[[1,1],[2,2]],[3,3]],[4,4]]", 445},
		{"[[[[3,0],[5,3]],[4,4]],[5,5]]", 791},
		{"[[[[5,0],[7,4]],[5,5]],[6,6]]", 1137},
		{"[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]", 3488},
		{"[9,1]", 29},
		{"[1,9]", 21},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := parseNumber(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			got, err := parsed.magnitude()

			if err != nil || got != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", got, err, tc.want)
			}
		})
	}
}

func TestFinalSums(t *testing.T) {
	tests := []struct {
		name, input string
		want        int
	}{
		{"four numbers", "[1,1]\n[2,2]\n[3,3]\n[4,4]", 445},
		{"five numbers", "[1,1]\n[2,2]\n[3,3]\n[4,4]\n[5,5]", 791},
		{"six numbers", "[1,1]\n[2,2]\n[3,3]\n[4,4]\n[5,5]\n[6,6]", 1137},
		{"slightly larger example", largerExample, 3488},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, _, err := solve(tc.input)

			if err != nil || part1 != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", part1, err, tc.want)
			}
		})
	}
}

func TestAdditionReducesTheSum(t *testing.T) {
	left, err := parseNumber("[[[[4,3],4],4],[7,[[8,4],9]]]")
	if err != nil {
		t.Fatal(err)
	}
	right, err := parseNumber("[1,1]")
	if err != nil {
		t.Fatal(err)
	}
	want, err := parseNumber("[[[[0,7],4],[[7,8],[6,0]]],[8,1]]")
	if err != nil {
		t.Fatal(err)
	}

	got := left.add(right)

	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSingleExplodes(t *testing.T) {
	tests := []struct{ before, after string }{
		{"[[[[[9,8],1],2],3],4]", "[[[[0,9],2],3],4]"},
		{"[7,[6,[5,[4,[3,2]]]]]", "[7,[6,[5,[7,0]]]]"},
		{"[[6,[5,[4,[3,2]]]],1]", "[[6,[5,[7,0]]],3]"},
		{"[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]"},
		{"[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[7,0]]]]"},
	}
	for _, tc := range tests {
		t.Run(tc.before, func(t *testing.T) {
			parsed, err := parseNumber(tc.before)
			if err != nil {
				t.Fatal(err)
			}
			want, err := parseNumber(tc.after)
			if err != nil {
				t.Fatal(err)
			}

			got, ok := parsed.explode()

			if !ok {
				t.Fatalf("%s should explode", tc.before)
			}
			if !slices.Equal(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func TestSplitRoundsDownThenUp(t *testing.T) {
	parsed, err := parseNumber("[11,0]")
	if err != nil {
		t.Fatal(err)
	}
	want, err := parseNumber("[[5,6],0]")
	if err != nil {
		t.Fatal(err)
	}

	got, ok := parsed.split()

	if !ok {
		t.Fatal("[11,0] should split")
	}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if _, ok := got.split(); ok {
		t.Fatal("nothing left to split")
	}
}

func TestSolveBoundaries(t *testing.T) {
	tests := []struct {
		name, input  string
		part1, part2 int
		wantError    bool
	}{
		// [[9,1],[1,9]] has magnitude 129, the reverse order only 121.
		{"addition order matters", "[9,1]\n[1,9]", 129, 129, false},
		{"windows newlines", "[9,1]\r\n[1,9]\r\n", 129, 129, false},
		{"single number", "[1,2]", 0, 0, true},
		{"empty input", "", 0, 0, true},
		{"blank input", "   \n\n", 0, 0, true},
		{"regular number only", "5\n[1,2]", 0, 0, true},
		{"unbalanced brackets", "[1,2\n[3,4]", 0, 0, true},
		{"trailing bracket", "[1,2]]\n[3,4]", 0, 0, true},
		{"trailing text", "[1,2]x\n[3,4]", 0, 0, true},
		{"triple", "[1,2,3]\n[3,4]", 0, 0, true},
		{"empty pair", "[]\n[1,2]", 0, 0, true},
		{"single element pair", "[1,2]\n[3,[4]]", 0, 0, true},
		{"negative number", "[1,-2]\n[3,4]", 0, 0, true},
		{"number too large", "[99999999999999999999,2]\n[3,4]", 0, 0, true},
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
