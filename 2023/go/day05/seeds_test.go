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

	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 35 {
		t.Fatalf("part1 = %d, want 35", part1)
	}
	if part2 != 46 {
		t.Fatalf("part2 = %d, want 46", part2)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing seeds", "\n\nseed-to-soil map:\n50 100 5"},
		{"missing map header", "seeds: 10 1\n\n50 100 5"},
		{"no maps", "seeds: 10 1"},
		{"seed not a number", "seeds: abc 1\n\nseed-to-soil map:\n50 100 5"},
		{"mapping not 3 fields", "seeds: 10 1\n\nseed-to-soil map:\n50 100"},
		{"mapping field not a number", "seeds: 10 1\n\nseed-to-soil map:\nabc 100 5"},
		{"empty seed list", "seeds: \n\nseed-to-soil map:\n50 100 5"},
		{"odd seed count", "seeds: 10 20 30\n\nseed-to-soil map:\n50 100 5"},
		{"empty mapping section", "seeds: 10 1\n\nseed-to-soil map:\n\nsoil-to-fertilizer map:\n1 1 1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := solve(tc.input)
			if err == nil {
				t.Fatalf("solve(%q) succeeded, want error", tc.input)
			}
		})
	}
}

func TestRangeApplication(t *testing.T) {
	tests := []struct {
		name       string
		range_     Range
		mapping    Mapping
		wantRanges []Range
	}{
		{
			"no overlap",
			Range{start: 10, length: 5},
			Mapping{dstStart: 100, srcStart: 20, length: 5},
			[]Range{{start: 10, length: 5}},
		},
		{
			"full overlap",
			Range{start: 20, length: 5},
			Mapping{dstStart: 100, srcStart: 20, length: 5},
			[]Range{{start: 100, length: 5}},
		},
		{
			"partial overlap",
			Range{start: 18, length: 5},
			Mapping{dstStart: 100, srcStart: 20, length: 5},
			[]Range{{start: 18, length: 2}, {start: 100, length: 3}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := AlmanacMap{
				name:     "test",
				mappings: []Mapping{tc.mapping},
			}
			got := m.ApplyRanges([]Range{tc.range_})

			if len(got) != len(tc.wantRanges) {
				t.Fatalf("got %d ranges, want %d: %v vs %v", len(got), len(tc.wantRanges), got, tc.wantRanges)
			}
			for i, r := range got {
				if r.start != tc.wantRanges[i].start || r.length != tc.wantRanges[i].length {
					t.Fatalf("range[%d] = {%d, %d}, want {%d, %d}",
						i, r.start, r.length, tc.wantRanges[i].start, tc.wantRanges[i].length)
				}
			}
		})
	}
}

func TestTransform(t *testing.T) {
	m := Mapping{dstStart: 100, srcStart: 50, length: 10}

	tests := []struct {
		src   int64
		dst   int64
		found bool
	}{
		{49, 49, false},
		{50, 100, true},
		{54, 104, true},
		{59, 109, true},
		{60, 60, false},
	}

	for _, tc := range tests {
		dst, found := m.Transform(tc.src)
		if found != tc.found || dst != tc.dst {
			t.Fatalf("Transform(%d) = (%d, %v), want (%d, %v)",
				tc.src, dst, found, tc.dst, tc.found)
		}
	}
}

func TestAlmanacApply(t *testing.T) {
	am := AlmanacMap{
		name: "test",
		mappings: []Mapping{
			{dstStart: 100, srcStart: 50, length: 10},
			{dstStart: 200, srcStart: 70, length: 5},
		},
	}

	tests := []struct {
		value int64
		want  int64
	}{
		{49, 49},  // before first range
		{50, 100}, // in first range
		{59, 109}, // in first range
		{65, 65},  // between ranges
		{70, 200}, // in second range
		{80, 80},  // after ranges
	}

	for _, tc := range tests {
		got := am.Apply(tc.value)
		if got != tc.want {
			t.Fatalf("Apply(%d) = %d, want %d", tc.value, got, tc.want)
		}
	}
}
