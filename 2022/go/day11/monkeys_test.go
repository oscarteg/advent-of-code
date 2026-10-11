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

	if err != nil || part1 != 10605 || part2 != 2713310158 {
		t.Fatalf("got (%d, %d, %v), want (10605, 2713310158, nil)", part1, part2, err)
	}
}

func TestTwoMonkeysSwapItemsEachRound(t *testing.T) {
	input := "Monkey 0:\n  Starting items: 1, 2\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 1\n    If false: throw to monkey 1\n\nMonkey 1:\n  Starting items: 3\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 0\n    If false: throw to monkey 0\n"
	monkeys, err := parse(input)
	if err != nil {
		t.Fatal(err)
	}
	answer := run(monkeys, 1, func(worry uint64) uint64 { return worry / reliefFactor })
	// Monkey 0 inspects 2 items, monkey 1 inspects the 3 items it now holds: 2 * 3 = 6.
	if answer != 6 {
		t.Fatalf("got %d, want 6", answer)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing fields", "Monkey 0:\n  Starting items: 1\n  Operation: new = old + 1\n"},
		{"bad operator", "Monkey 0:\n  Starting items: 1\n  Operation: new = old ^ 1\n  Test: divisible by 2\n    If true: throw to monkey 0\n    If false: throw to monkey 0\n"},
		{"throw target out of range", "Monkey 0:\n  Starting items: 1\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 1\n    If false: throw to monkey 0\n"},
		{"bad number", "Monkey 0:\n  Starting items: x\n  Operation: new = old + 1\n  Test: divisible by 2\n    If true: throw to monkey 0\n    If false: throw to monkey 0\n"},
		{"empty input", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("expected error for input %q", tc.input)
			}
		})
	}
}
