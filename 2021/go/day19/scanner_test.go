package main

import (
	"os"
	"testing"
)

// sharedTwelve lists the twelve beacons scanners 0 and 1 share, as the puzzle
// writes them in each scanner's own frame; scanner 1 sits at 68,-1246,-43.
const sharedTwelve = "--- scanner 0 ---\n" +
	"-618,-824,-621\n-537,-823,-458\n-447,-329,318\n404,-588,-901\n" +
	"544,-627,-890\n528,-643,409\n-661,-816,-575\n390,-675,-793\n" +
	"423,-701,434\n-345,-311,381\n459,-707,401\n-485,-357,347\n\n" +
	"--- scanner 1 ---\n" +
	"686,422,578\n605,423,415\n515,917,-361\n-336,658,858\n" +
	"-476,619,847\n-460,603,-452\n729,430,532\n-322,571,750\n" +
	"-355,545,-477\n413,935,-424\n-391,539,-444\n553,889,-390\n"

// orientations lists the puzzle's six beacons in five of the 24 orientations.
const orientations = "--- scanner 0 ---\n-1,-1,1\n-2,-2,2\n-3,-3,3\n-2,-3,1\n5,6,-4\n8,0,7\n\n" +
	"--- scanner 0 ---\n1,-1,1\n2,-2,2\n3,-3,3\n2,-1,3\n-5,4,-6\n-8,-7,0\n\n" +
	"--- scanner 0 ---\n-1,-1,-1\n-2,-2,-2\n-3,-3,-3\n-1,-3,-2\n4,6,5\n-7,0,8\n\n" +
	"--- scanner 0 ---\n1,1,-1\n2,2,-2\n3,3,-3\n1,3,-2\n-4,-6,5\n7,0,8\n\n" +
	"--- scanner 0 ---\n1,1,1\n2,2,2\n3,3,3\n3,1,2\n-6,-4,-5\n0,7,-8\n"

func TestSolve(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		input        string
		beacons      int
		largestSpace int
	}{
		{"official example", string(example), 79, 3621},
		{"twelve shared beacons", sharedTwelve, 12, 1357},
		{"single report needs no alignment", "--- scanner 0 ---\n1,2,3\n4,5,6\n", 2, 0},
		{"repeated beacon counts once", "--- scanner 0 ---\n1,2,3\n1,2,3\n", 1, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			beacons, largestSpace, err := solve(tc.input)

			if err != nil || beacons != tc.beacons || largestSpace != tc.largestSpace {
				t.Fatalf("got (%d, %d, %v), want (%d, %d, nil)", beacons, largestSpace, err, tc.beacons, tc.largestSpace)
			}
		})
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty input", ""},
		{"beacon before header", "1,2,3\n"},
		{"two coordinates", "--- scanner 0 ---\n1,2\n"},
		{"non numeric coordinate", "--- scanner 0 ---\n1,2,x\n"},
		{"non numeric scanner id", "--- scanner one ---\n1,2,3\n"},
		{"reports without overlap", "--- scanner 0 ---\n1,2,3\n\n--- scanner 1 ---\n4,5,6\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want an error", tc.input)
			}
		})
	}
}

func TestRotationsReorientOneBeaconSet(t *testing.T) {
	rotations := allRotations()
	if len(rotations) != 24 {
		t.Fatalf("got %d rotations, want 24", len(rotations))
	}
	reports, err := parseReports(orientations)
	if err != nil {
		t.Fatal(err)
	}
	reference := make(map[point]bool, len(reports[0].beacons))
	for _, beacon := range reports[0].beacons {
		reference[beacon] = true
	}
	for index, r := range reports[1:] {
		if !matchesUnderSomeRotation(r, reference, rotations) {
			t.Fatalf("orientation %d matches no rotation of the first report", index+1)
		}
	}
}

func matchesUnderSomeRotation(r report, reference map[point]bool, rotations []rotation) bool {
	for _, rot := range rotations {
		matched := 0
		for _, beacon := range r.beacons {
			if reference[rot.apply(beacon)] {
				matched++
			}
		}
		if matched == len(reference) {
			return true
		}
	}
	return false
}
