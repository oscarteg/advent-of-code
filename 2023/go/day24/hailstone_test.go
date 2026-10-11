package main

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
)

const (
	exampleAreaLow  = 7
	exampleAreaHigh = 27
)

const threeHailstones = "12, 22, 31 @ -1, 0, 2\n" +
	"8, 30, 34 @ 2, -3, 1\n" +
	"13, 11, 45 @ 0, 5, -2\n"

const realAreaCrossing = "200000000000000, 200000000000000, 0 @ 1000, 1000, 0\n" +
	"400000000000000, 200000000000000, 0 @ -1000, 1000, 0\n"

func readExample(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	return string(input)
}

func TestOfficialExample(t *testing.T) {
	crossings, rock, err := solveWith(readExample(t), exampleAreaLow, exampleAreaHigh)

	if err != nil || crossings != 2 || rock != 47 {
		t.Fatalf("got (%d, %d, %v), want (2, 47, nil)", crossings, rock, err)
	}
}

func TestCountCrossings(t *testing.T) {
	tests := []struct {
		name, input string
		low, high   int64
		want        int64
	}{
		{"official example", readExample(t), exampleAreaLow, exampleAreaHigh, 2},
		{"crossing on the corner counts", "0, 0, 0 @ 1, 1, 0\n4, 0, 0 @ 0, 1, 0", 0, 4, 1},
		{"crossing outside the area", "0, 0, 0 @ 1, 1, 0\n4, 0, 0 @ 0, 1, 0", 0, 3, 0},
		{"parallel paths never cross", "0, 0, 0 @ 1, 1, 0\n1, 0, 0 @ 2, 2, 0", -10, 10, 0},
		{"crossed in the past", "0, 0, 0 @ -1, -1, 0\n4, 0, 0 @ 0, 1, 0", -10, 10, 0},
		{"the z axis is ignored", "0, 0, 0 @ 1, 1, 9\n4, 0, 0 @ 0, 1, -9", 0, 4, 1},
		{"the real test area stays exact", realAreaCrossing, testAreaLow, testAreaHigh, 1},
		{"a crossing beyond the area is excluded", realAreaCrossing, testAreaLow, 299_999_999_999_999, 0},
		{"carriage returns and spacing", "19, 13, 30 @ -2,  1, -2\r\n18, 19, 22 @ -1, -1, -2\r\n", exampleAreaLow, exampleAreaHigh, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stones, err := parseHailstones(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if got := countCrossings(stones, tc.low, tc.high); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestRockPositionSum(t *testing.T) {
	example := readExample(t)
	reversed := strings.Split(strings.TrimSpace(example), "\n")
	slices.Reverse(reversed)

	tests := []struct {
		name, input string
		want        int64
	}{
		{"official example", example, 47},
		{"three hailstones determine the rock", threeHailstones, 60},
		{"extra hailstones agree", threeHailstones + "-18, 27, 58 @ 5, 1, -1\n", 60},
		{"the order does not matter", strings.Join(reversed, "\n"), 47},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stones, err := parseHailstones(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			got, err := rockPositionSum(stones)

			if err != nil || got != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", got, err, tc.want)
			}
		})
	}
}

func TestInvalidInput(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty input", ""},
		{"missing separator", "19, 13, 30"},
		{"missing position coordinate", "19, 13 @ -2, 1, -2"},
		{"missing velocity coordinate", "19, 13, 30 @ -2, 1"},
		{"coordinate is not a number", "x, 13, 30 @ -2, 1, -2"},
		{"too many coordinates", "19, 13, 30, 7 @ -2, 1, -2, 3"},
		{"one hailstone cannot determine a rock", "19, 13, 30 @ -2, 1, -2"},
		{"two hailstones cannot determine a rock", "19, 13, 30 @ -2, 1, -2\n18, 19, 22 @ -1, -1, -2"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := solve(tc.input)

			if !errors.Is(err, errInvalidHailstones) {
				t.Fatalf("got %v, want an errInvalidHailstones error", err)
			}
		})
	}
}
