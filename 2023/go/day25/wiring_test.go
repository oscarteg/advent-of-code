package main

import (
	"os"
	"testing"
)

func exampleInput(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	return string(input)
}

func TestOfficialExample(t *testing.T) {
	part1, err := solve(exampleInput(t))

	if err != nil || part1 != 54 {
		t.Fatalf("got (%d, %v), want (54, nil)", part1, err)
	}
}

func TestExampleGroupSizes(t *testing.T) {
	d, err := parseDiagram(exampleInput(t))
	if err != nil {
		t.Fatal(err)
	}
	if d.components() != 15 {
		t.Fatalf("got %d components, want 15", d.components())
	}

	left, right := d.minimumCut()

	if left > right {
		left, right = right, left
	}
	if left != 6 || right != 9 {
		t.Fatalf("got groups (%d, %d), want (6, 9)", left, right)
	}
}

func TestSmallDiagrams(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "two components on a single wire", input: "jqt: rhn", want: 1},
		{name: "triangle", input: "a: b c\nb: c", want: 2},
		{name: "triangle wired in the other direction", input: "c: a b\nb: a", want: 2},
		{name: "triangles joined by one wire", input: "a: b c\nb: c\nc: d\nd: e f\ne: f", want: 9},
		{
			name:  "five-component cliques joined by three wires",
			input: "a: b c d e f\nb: c d e g\nc: d e h\nd: e\nf: g h i j\ng: h i j\nh: i j\ni: j",
			want:  25,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := solve(tc.input)
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMalformedInputReturnsErrors(t *testing.T) {
	for _, input := range []string{
		"",
		"jqt rhn xhk",
		"jqt:",
		": rhn",
		"jqt: jqt",
	} {
		if _, err := solve(input); err == nil {
			t.Fatalf("expected error for input %q", input)
		}
	}
}
