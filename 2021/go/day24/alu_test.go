package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func blockText(parameters block) string {
	return strings.Join([]string{
		"inp w", "mul x 0", "add x z", "mod x 26",
		"div z " + strconv.FormatInt(parameters.div, 10),
		"add x " + strconv.FormatInt(parameters.addX, 10),
		"eql x w", "eql x 0", "mul y 0", "add y 25", "mul y x", "add y 1",
		"mul z y", "mul y 0", "add y w",
		"add y " + strconv.FormatInt(parameters.addY, 10),
		"mul y x", "add z y", "",
	}, "\n")
}

func programText(blocks ...block) string {
	var program strings.Builder
	for _, parameters := range blocks {
		program.WriteString(blockText(parameters))
	}
	return program.String()
}

func digitsOf(number uint64) ([]int64, bool) {
	text := strconv.FormatUint(number, 10)
	digits := make([]int64, 0, len(text))
	for _, character := range text {
		if character == '0' {
			return nil, false
		}
		digits = append(digits, int64(character-'0'))
	}
	return digits, true
}

func TestRunOfficialPrograms(t *testing.T) {
	const binarySplit = "inp w\nadd z w\nmod z 2\ndiv w 2\nadd y w\nmod y 2\ndiv w 2\nadd x w\nmod x 2\ndiv w 2\nmod w 2\n"
	const threeTimes = "inp z\ninp x\nmul z 3\neql z x\n"
	tests := []struct {
		name    string
		program string
		inputs  []int64
		want    [4]int64
	}{
		{name: "negation", program: "inp x\nmul x -1\n", inputs: []int64{7}, want: [4]int64{0, -7, 0, 0}},
		{name: "three times larger", program: threeTimes, inputs: []int64{3, 9}, want: [4]int64{0, 9, 0, 1}},
		{name: "not three times larger", program: threeTimes, inputs: []int64{3, 8}, want: [4]int64{0, 8, 0, 0}},
		{name: "binary split of 11", program: binarySplit, inputs: []int64{11}, want: [4]int64{1, 0, 1, 1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, err := parseProgram(test.program)
			if err != nil {
				t.Fatalf("parseProgram() error = %v", err)
			}
			registers, err := run(program, test.inputs)
			if err != nil {
				t.Fatalf("run() error = %v", err)
			}
			if registers != test.want {
				t.Errorf("run() = %v, want %v", registers, test.want)
			}
		})
	}
}

func TestSolve(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	tests := []struct {
		name         string
		input        string
		part1, part2 uint64
	}{
		{name: "example monad", input: string(example), part1: 99999596891916, part2: 49411151121911},
		{
			name:  "pair eight apart forces both digits",
			input: programText(block{div: 1, addX: 13}, block{div: 26, addX: 8}),
			part1: 19,
			part2: 19,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			part1, part2, err := solve(test.input)
			if err != nil {
				t.Fatalf("solve() error = %v", err)
			}
			if part1 != test.part1 || part2 != test.part2 {
				t.Errorf("solve() = %d, %d, want %d, %d", part1, part2, test.part1, test.part2)
			}
		})
	}
}

func TestSolveMatchesBruteForce(t *testing.T) {
	input := programText(
		block{div: 1, addX: 13, addY: 6},
		block{div: 1, addX: 11, addY: 9},
		block{div: 26, addX: -4, addY: 7},
		block{div: 26, addX: -8, addY: 5},
	)
	parsed, err := parseMonad(input)
	if err != nil {
		t.Fatalf("parseMonad() error = %v", err)
	}
	var smallestValid, largestValid uint64
	for number := uint64(1111); number <= 9999; number++ {
		digits, ok := digitsOf(number)
		if !ok {
			continue
		}
		registers, err := run(parsed.program, digits)
		if err != nil {
			t.Fatalf("run(%d) error = %v", number, err)
		}
		if registers[registerZ] != 0 {
			continue
		}
		if smallestValid == 0 {
			smallestValid = number
		}
		largestValid = number
	}
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve() error = %v", err)
	}
	if part1 != largestValid || part2 != smallestValid {
		t.Errorf("solve() = %d, %d, want %d, %d", part1, part2, largestValid, smallestValid)
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	twoBlocks := programText(block{div: 1, addX: 13, addY: 6}, block{div: 26, addX: -4, addY: 7})
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty program", input: ""},
		{name: "unknown register", input: "inp q\n"},
		{name: "unknown instruction", input: "nop x 1\n"},
		{name: "missing target", input: "inp\n"},
		{name: "missing operand", input: "add x\n"},
		{name: "trailing tokens", input: "mul x y z\n"},
		{name: "partial block", input: "inp w\nmul x 0\n"},
		{name: "not the monad template", input: strings.Replace(twoBlocks, "mod x 26", "mod x 25", 1)},
		{name: "never popped", input: programText(block{div: 1, addX: 13, addY: 6}, block{div: 1, addX: 11, addY: 9})},
		{name: "pop before push", input: programText(block{div: 26, addX: -4, addY: 7}, block{div: 1, addX: 13, addY: 6})},
		{name: "unsatisfiable pair", input: programText(block{div: 1, addX: 13}, block{div: 26, addX: 9})},
		{name: "push that can match a digit", input: programText(block{div: 1, addX: 5}, block{div: 26, addX: -4})},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := solve(test.input); err == nil {
				t.Error("solve() error = nil, want an error")
			}
		})
	}
}

func TestRunRejectsCrashingPrograms(t *testing.T) {
	tests := []struct {
		name    string
		program string
		inputs  []int64
	}{
		{name: "div by zero", program: "inp x\ndiv x 0\n", inputs: []int64{1}},
		{name: "mod of a negative value", program: "inp x\nmul x -1\nmod x 2\n", inputs: []int64{1}},
		{name: "missing input", program: "inp x\n", inputs: nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, err := parseProgram(test.program)
			if err != nil {
				t.Fatalf("parseProgram() error = %v", err)
			}
			if _, err := run(program, test.inputs); err == nil {
				t.Error("run() error = nil, want an error")
			}
		})
	}
}
