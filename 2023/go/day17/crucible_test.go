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

	if err != nil || part1 != 102 || part2 != 94 {
		t.Fatalf("got (%d, %d, %v), want (102, 94, nil)", part1, part2, err)
	}
}

func TestUltraCrucibleStraightCorridor(t *testing.T) {
	const input = "111111111111\n999999999991\n999999999991\n999999999991\n999999999991\n"

	g, err := parseGrid(input)
	if err != nil {
		t.Fatal(err)
	}
	heatLoss, err := leastHeatLoss(g, 4, 10)

	if err != nil || heatLoss != 71 {
		t.Fatalf("got (%d, %v), want (71, nil)", heatLoss, err)
	}
}

func TestSingleBlockGridNeedsNoMoves(t *testing.T) {
	g, err := parseGrid("5")
	if err != nil {
		t.Fatal(err)
	}

	heatLoss, err := leastHeatLoss(g, 0, 3)

	if err != nil || heatLoss != 0 {
		t.Fatalf("got (%d, %v), want (0, nil)", heatLoss, err)
	}
}

func TestStraightLineLongerThanThreeMustDetour(t *testing.T) {
	g, err := parseGrid("1111\n9991\n1991\n1119")
	if err != nil {
		t.Fatal(err)
	}

	heatLoss, err := leastHeatLoss(g, 0, 3)

	if err != nil || heatLoss != 14 {
		t.Fatalf("got (%d, %v), want (14, nil)", heatLoss, err)
	}
}

func TestGridTooSmallForMinimumRunHasNoPath(t *testing.T) {
	_, _, err := solve("12\n34")

	if err == nil {
		t.Fatal("expected an error for a grid too small to satisfy the minimum run")
	}
}

func TestMalformedInputIsRejected(t *testing.T) {
	for _, input := range []string{"", "12\n1", "1a2\n345"} {
		if _, _, err := solve(input); err == nil {
			t.Errorf("solve(%q): expected an error", input)
		}
	}
}
