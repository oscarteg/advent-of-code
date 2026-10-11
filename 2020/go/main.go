package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/oscarteg/advent-of-code-2020/challenge/day1"
	"github.com/oscarteg/advent-of-code-2020/challenge/day10"
	"github.com/oscarteg/advent-of-code-2020/challenge/day11"
	"github.com/oscarteg/advent-of-code-2020/challenge/day12"
	"github.com/oscarteg/advent-of-code-2020/challenge/day13"
	"github.com/oscarteg/advent-of-code-2020/challenge/day14"
	"github.com/oscarteg/advent-of-code-2020/challenge/day15"
	"github.com/oscarteg/advent-of-code-2020/challenge/day16"
	"github.com/oscarteg/advent-of-code-2020/challenge/day17"
	"github.com/oscarteg/advent-of-code-2020/challenge/day18"
	"github.com/oscarteg/advent-of-code-2020/challenge/day19"
	"github.com/oscarteg/advent-of-code-2020/challenge/day2"
	"github.com/oscarteg/advent-of-code-2020/challenge/day20"
	"github.com/oscarteg/advent-of-code-2020/challenge/day21"
	"github.com/oscarteg/advent-of-code-2020/challenge/day22"
	"github.com/oscarteg/advent-of-code-2020/challenge/day23"
	"github.com/oscarteg/advent-of-code-2020/challenge/day24"
	"github.com/oscarteg/advent-of-code-2020/challenge/day25"
	"github.com/oscarteg/advent-of-code-2020/challenge/day3"
	"github.com/oscarteg/advent-of-code-2020/challenge/day4"
	"github.com/oscarteg/advent-of-code-2020/challenge/day5"
	"github.com/oscarteg/advent-of-code-2020/challenge/day6"
	"github.com/oscarteg/advent-of-code-2020/challenge/day7"
	"github.com/oscarteg/advent-of-code-2020/challenge/day8"
	"github.com/oscarteg/advent-of-code-2020/challenge/day9"
)

type runner func(input string, args []string) error

// both runs the two solvers of a day on the same stdin input.
func both[A, B any](partA func(*challenge.Input) (A, error), partB func(*challenge.Input) (B, error)) runner {
	return func(input string, _ []string) error {
		answer, err := partA(challenge.FromLiteral(input))
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
		fmt.Printf("Part 1: %v\n", answer)
		answerB, err := partB(challenge.FromLiteral(input))
		if err != nil {
			return fmt.Errorf("part 2: %w", err)
		}
		fmt.Printf("Part 2: %v\n", answerB)
		return nil
	}
}

var days = map[string]runner{
	"1":  both(day1.A, day1.B),
	"2":  both(day2.A, day2.B),
	"3":  both(day3.A, day3.B),
	"4":  both(day4.A, day4.B),
	"5":  both(day5.A, day5.B),
	"6":  both(day6.A, day6.B),
	"7":  both(day7.A, day7.B),
	"8":  both(day8.A, day8.B),
	"10": both(day10.A, day10.B),
	"11": both(day11.A, day11.B),
	"12": both(day12.A, day12.B),
	"13": both(day13.A, day13.B),
	"14": both(day14.A, day14.B),
	"15": both(day15.A, day15.B),
	"16": both(day16.A, day16.B),
	"17": both(day17.A, day17.B),
	"18": both(day18.A, day18.B),
	"19": both(day19.A, day19.B),
	"20": both(day20.A, day20.B),
	"21": both(day21.A, day21.B),
	"22": both(day22.A, day22.B),
	"23": both(day23.A, day23.B),
	"24": both(day24.A, day24.B),
	// Day 9 takes an optional preamble length (the official example uses 5).
	"9": func(input string, args []string) error {
		preamble := day9.RealPreamble
		if len(args) > 0 {
			p, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("preamble %q: %w", args[0], err)
			}
			preamble = p
		}
		return both(
			func(in *challenge.Input) (int64, error) { return day9.A(in, preamble) },
			func(in *challenge.Input) (int64, error) { return day9.B(in, preamble) },
		)(input, nil)
	},
	// Day 25 has no part 2.
	"25": func(input string, _ []string) error {
		answer, err := day25.A(challenge.FromLiteral(input))
		if err != nil {
			return fmt.Errorf("part 1: %w", err)
		}
		fmt.Printf("Part 1: %v\n", answer)
		return nil
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run . <day> [day 9: preamble] < input.txt")
		os.Exit(2)
	}
	run, ok := days[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown day %q\n", os.Args[1])
		os.Exit(2)
	}
	input, err := io.ReadAll(os.Stdin)
	if err == nil {
		err = run(string(input), os.Args[2:])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
