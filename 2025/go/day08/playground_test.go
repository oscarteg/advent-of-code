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

	if err != nil || part1 != 40 || part2 != 25272 {
		t.Fatalf("got (%d, %d, %v), want (40, 25272, nil)", part1, part2, err)
	}
}

func TestConnectingCircuits(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		connectionCount int
		part1, part2    int
		wantError       bool
	}{
		{
			name:            "two close pairs bridged by the long connection between them",
			input:           "0,0,0\n1,0,0\n10,0,0\n11,0,0\n",
			connectionCount: 1,
			part1:           2,
			part2:           10,
		},
		{
			name:            "a pair already sharing a circuit counts but changes nothing",
			input:           "0,0,0\n1,0,0\n2,0,0\n100,0,0\n200,0,0\n",
			connectionCount: 3,
			part1:           3,
			part2:           20000,
		},
		{
			name:            "surrounding whitespace is ignored",
			input:           "  0,0,0 \n\n 1,0,0\n10,0,0\n 11,0,0 \n\n",
			connectionCount: 1,
			part1:           2,
			part2:           10,
		},
		{
			name:            "fewer than three circuits left",
			input:           "0,0,0\n1,0,0\n",
			connectionCount: 1,
			wantError:       true,
		},
		{
			name:      "box missing a coordinate",
			input:     "1,2\n3,4,5\n",
			wantError: true,
		},
		{
			name:      "box with a fourth coordinate",
			input:     "1,2,3,4\n5,6,7\n",
			wantError: true,
		},
		{
			name:      "non-numeric coordinate",
			input:     "1,2,x\n3,4,5\n",
			wantError: true,
		},
		{
			name:      "single box",
			input:     "1,2,3\n",
			wantError: true,
		},
		{
			name:      "empty input",
			input:     "  \n \n",
			wantError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solveWith(tc.input, tc.connectionCount)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
