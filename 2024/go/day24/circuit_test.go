package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const smallExample = `x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02
`

type wireSwap struct{ left, right string }

func wireName(prefix byte, bit int) string {
	return fmt.Sprintf("%c%02d", prefix, bit)
}

func swapName(name string, swaps []wireSwap) string {
	for _, swap := range swaps {
		switch name {
		case swap.left:
			return swap.right
		case swap.right:
			return swap.left
		}
	}
	return name
}

// Builds a correct ripple-carry adder for x + y, then renames the gate outputs listed in
// swaps so the circuit is wrong in exactly the way part 2 describes.
func adderInput(bits int, x, y int64, swaps []wireSwap) string {
	var text strings.Builder
	for _, side := range []struct {
		prefix byte
		value  int64
	}{{'x', x}, {'y', y}} {
		for bit := range bits {
			fmt.Fprintf(&text, "%s: %d\n", wireName(side.prefix, bit), (side.value>>bit)&1)
		}
	}
	text.WriteByte('\n')

	gate := func(left, op, right, output string) {
		fmt.Fprintf(&text, "%s %s %s -> %s\n", left, op, right, swapName(output, swaps))
	}
	firstCarry := wireName('c', 0)
	if bits == 1 {
		firstCarry = wireName('z', 1)
	}
	gate(wireName('x', 0), "XOR", wireName('y', 0), wireName('z', 0))
	gate(wireName('x', 0), "AND", wireName('y', 0), firstCarry)

	for bit := 1; bit < bits; bit++ {
		carryIn, sum, generated := wireName('c', bit-1), wireName('s', bit), wireName('a', bit)
		carryOut := wireName('c', bit)
		if bit+1 == bits {
			carryOut = wireName('z', bits)
		}
		gate(wireName('x', bit), "XOR", wireName('y', bit), sum)
		gate(wireName('x', bit), "AND", wireName('y', bit), generated)
		gate(sum, "XOR", carryIn, wireName('z', bit))
		gate(sum, "AND", carryIn, wireName('p', bit))
		gate(wireName('p', bit), "OR", generated, carryOut)
	}
	return text.String()
}

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, _, err := solve(string(input))
	if err != nil || part1 != 2024 {
		t.Fatalf("got (%d, %v), want (2024, nil)", part1, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		part1 int64
		part2 string
	}{
		{"small official example", smallExample, 4, ""},
		{
			"a correct adder adds its inputs and needs no swaps",
			adderInput(8, 0b1011_0101, 0b0110_1110, nil),
			0b1011_0101 + 0b0110_1110,
			"",
		},
		{"edge case: a half adder is the smallest correct adder", adderInput(1, 1, 1, nil), 2, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve() error = %v", err)
			}
			if part1 != tc.part1 {
				t.Errorf("part1 = %d, want %d", part1, tc.part1)
			}
			if part2 != tc.part2 {
				t.Errorf("part2 = %q, want %q", part2, tc.part2)
			}
		})
	}
}

func TestSolveFindsCrossedAdderWires(t *testing.T) {
	tests := []struct {
		name  string
		swaps []wireSwap
		want  string
	}{
		{"two crossed pairs", []wireSwap{{"z02", "c02"}, {"s05", "a05"}}, "a05,c02,s05,z02"},
		{
			"four crossed pairs",
			[]wireSwap{{"s01", "a01"}, {"z03", "a03"}, {"z05", "c05"}, {"s06", "a06"}},
			"a01,a03,a06,c05,s01,s06,z03,z05",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, part2, err := solve(adderInput(8, 0b1011_0101, 0b0110_1110, tc.swaps))
			if err != nil {
				t.Fatalf("solve() error = %v", err)
			}
			if part2 != tc.want {
				t.Errorf("part2 = %q, want %q", part2, tc.want)
			}
		})
	}
}

func TestSolveMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"no gates", "x00: 1\ny00: 1\n"},
		{"initial value is not a bit", "x00: 2\ny00: 1\n\nx00 AND y00 -> z00\n"},
		{"unknown operator", "x00: 1\ny00: 1\n\nx00 NAND y00 -> z00\n"},
		{"line is neither a value nor a gate", "x00 1\n\nx00 AND y00 -> z00\n"},
		{"gate is missing an input", "x00: 1\ny00: 1\n\nx00 XOR -> z00\n"},
		{"input wire never receives a value", "x00: 1\n\nx00 AND y00 -> z00\n"},
		{"wire driven twice", "x00: 1\ny00: 1\n\nx00 AND y00 -> z00\nx00 OR y00 -> z00\n"},
		{"circuit has a cycle", "y00: 1\n\nz01 AND y00 -> z00\nz00 OR y00 -> z01\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatal("solve() error = nil, want an error")
			}
		})
	}
}
