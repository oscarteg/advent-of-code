package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 6440 {
		t.Errorf("part1 = %d, want 6440", part1)
	}

	if part2 != 5905 {
		t.Errorf("part2 = %d, want 5905", part2)
	}
}

func TestParseSingleHand(t *testing.T) {
	hand, err := parseHand("32T3K 765")
	if err != nil {
		t.Fatalf("parseHand failed: %v", err)
	}

	if hand.bid != 765 {
		t.Errorf("bid = %d, want 765", hand.bid)
	}

	if hand.handType != OnePair {
		t.Errorf("handType = %v, want OnePair", hand.handType)
	}
}

func TestParseInvalidHand(t *testing.T) {
	_, err := parseHand("INVALID 100")
	if err == nil {
		t.Error("parseHand should fail for invalid card")
	}
}

func TestClassifyHands(t *testing.T) {
	tests := []struct {
		line     string
		expected HandType
	}{
		{"AAAAA 1", FiveOfAKind},
		{"AA8AA 1", FourOfAKind},
		{"23332 1", FullHouse},
		{"TTT98 1", ThreeOfAKind},
		{"23432 1", TwoPair},
		{"A23A4 1", OnePair},
		{"23456 1", HighCard},
	}

	for _, tt := range tests {
		hand, err := parseHand(tt.line)
		if err != nil {
			t.Fatalf("parseHand failed: %v", err)
		}

		if hand.handType != tt.expected {
			t.Errorf("classifyHand(%s) = %v, want %v", tt.line, hand.handType, tt.expected)
		}
	}
}

func TestParseHandJokerAllJokers(t *testing.T) {
	hand, err := parseHandJoker("JJJJJ 100")
	if err != nil {
		t.Fatalf("parseHandJoker failed: %v", err)
	}

	if hand.handType != FiveOfAKind {
		t.Errorf("handType = %v, want FiveOfAKind", hand.handType)
	}
}
