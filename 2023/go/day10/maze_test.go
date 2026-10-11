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

	if err != nil || part1 != 8 || part2 != 1 {
		t.Fatalf("got (%d, %d, %v), want (8, 1, nil)", part1, part2, err)
	}
}

func TestEnclosedTileExamples(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name: "simple square loop",
			input: "" +
				".....\n" +
				".S-7.\n" +
				".|.|.\n" +
				".L-J.\n" +
				".....",
			want: 1,
		},
		{
			name: "large loop with junk pipes",
			input: "" +
				"...........\n" +
				".S-------7.\n" +
				".|F-----7|.\n" +
				".||.....||.\n" +
				".||.....||.\n" +
				".|L-7.F-J|.\n" +
				".|..|.|..|.\n" +
				".L--J.L--J.\n" +
				"...........",
			want: 4,
		},
		{
			name: "squeezed loop",
			input: "" +
				".F----7F7F7F7F-7....\n" +
				".|F--7||||||||FJ....\n" +
				".||.FJ||||||||L7....\n" +
				"FJL7L7LJLJ||LJ.L-7..\n" +
				"L--J.L7...LJS7F-7L7.\n" +
				"....F-J..F7FJ|L7L7L7\n" +
				"....L7.F7||L7|.L7L7|\n" +
				".....|FJLJ|FJ|F7|.LJ\n" +
				"....FJL-7.||.||||...\n" +
				"....L---J.LJ.LJLJ...",
			want: 8,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if part2 != tc.want {
				t.Fatalf("got %d, want %d", part2, tc.want)
			}
		})
	}
}

func TestSmallestPossibleLoop(t *testing.T) {
	part1, part2, err := solve("S7\nLJ")

	if err != nil || part1 != 2 || part2 != 0 {
		t.Fatalf("got (%d, %d, %v), want (2, 0, nil)", part1, part2, err)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []string{
		"",
		".....\n.X-7.\n.|.|.\n.L-J.\n.....",
		".....\n.F-7.\n.|.|\n.L-J.\n.....",
		".....\n..S..\n.....",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, _, err := solve(input); err == nil {
				t.Fatalf("expected an error for input %q", input)
			}
		})
	}
}
