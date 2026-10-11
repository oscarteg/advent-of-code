package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Mapping represents a range transformation in an almanac map.
type Mapping struct {
	dstStart int64
	srcStart int64
	length   int64
}

// Transform applies the mapping to a source value, returning (destination, true)
// if the value is in the source range, or (value, false) if unmapped.
func (m Mapping) Transform(src int64) (int64, bool) {
	if src >= m.srcStart && src < m.srcStart+m.length {
		return m.dstStart + (src - m.srcStart), true
	}
	return src, false
}

// AlmanacMap represents a stage in the transformation pipeline.
type AlmanacMap struct {
	name     string
	mappings []Mapping
}

// Apply maps a single value through this stage.
func (am AlmanacMap) Apply(value int64) int64 {
	for _, m := range am.mappings {
		if dst, ok := m.Transform(value); ok {
			return dst
		}
	}
	return value
}

// ApplyRanges maps a set of ranges through this stage, splitting at boundaries
// and transforming overlaps. Returns the resulting ranges sorted by start.
func (am AlmanacMap) ApplyRanges(ranges []Range) []Range {
	var result []Range

	for _, r := range ranges {
		rangeStart := r.start
		rangeEnd := r.start + r.length
		var current int64 = rangeStart

		// Process the range segment by segment
		for current < rangeEnd {
			// Find the earliest boundary after current
			nextBoundary := rangeEnd

			// Check all mapping boundaries
			for _, m := range am.mappings {
				// Mapping covers [m.srcStart, m.srcStart + m.length)
				if m.srcStart > current && m.srcStart < nextBoundary {
					nextBoundary = m.srcStart
				}
				mapEnd := m.srcStart + m.length
				if mapEnd > current && mapEnd < nextBoundary {
					nextBoundary = mapEnd
				}
			}

			// Now process [current, nextBoundary)
			segStart := current
			segEnd := nextBoundary

			// Check which mapping (if any) covers this segment
			var mapped bool
			for _, m := range am.mappings {
				mapStart := m.srcStart
				mapEnd := m.srcStart + m.length

				// Check if [segStart, segEnd) overlaps with [mapStart, mapEnd)
				overlapStart := max(segStart, mapStart)
				overlapEnd := min(segEnd, mapEnd)

				if overlapStart < overlapEnd {
					// This segment is within a mapping, transform it
					newStart := m.dstStart + (overlapStart - mapStart)
					newEnd := m.dstStart + (overlapEnd - mapStart)
					result = append(result, Range{start: newStart, length: newEnd - newStart})
					mapped = true
					break
				}
			}

			if !mapped {
				// This segment is not mapped, keep it as-is
				result = append(result, Range{start: segStart, length: segEnd - segStart})
			}

			current = nextBoundary
		}
	}

	// Sort by start position for consistent output
	slices.SortFunc(result, func(a, b Range) int {
		if a.start < b.start {
			return -1
		}
		if a.start > b.start {
			return 1
		}
		return 0
	})

	return result
}

// Range represents a half-open interval [start, start+length).
type Range struct {
	start  int64
	length int64
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Almanac holds the seeds and transformation pipeline.
type Almanac struct {
	seeds []int64
	maps  []AlmanacMap
}

// LocationOf traces a seed through all maps to get its final location.
func (a Almanac) LocationOf(seed int64) int64 {
	value := seed
	for _, m := range a.maps {
		value = m.Apply(value)
	}
	return value
}

// LocationOfRanges traces seed ranges through all maps and returns the minimum
// location value among the resulting ranges.
func (a Almanac) LocationOfRanges(seedRanges []Range) int64 {
	ranges := seedRanges
	for _, m := range a.maps {
		ranges = m.ApplyRanges(ranges)
	}

	if len(ranges) == 0 {
		return 0
	}

	minLoc := ranges[0].start
	for _, r := range ranges[1:] {
		if r.start < minLoc {
			minLoc = r.start
		}
	}
	return minLoc
}

// Parse reads the input and constructs an Almanac.
func ParseAlmanac(input string) (Almanac, error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	var a Almanac
	var i int

	// Parse seeds line
	if i >= len(lines) {
		return a, fmt.Errorf("parse almanac: missing seeds line")
	}
	seedLine := strings.TrimSpace(lines[i])
	if !strings.HasPrefix(seedLine, "seeds:") {
		return a, fmt.Errorf("parse almanac: first line must start with 'seeds:', got %q", seedLine)
	}
	seedsStr := strings.TrimPrefix(seedLine, "seeds:")
	seedStrs := strings.Fields(strings.TrimSpace(seedsStr))
	if len(seedStrs) == 0 {
		return a, fmt.Errorf("parse almanac: no seeds found")
	}
	a.seeds = make([]int64, len(seedStrs))
	for j, s := range seedStrs {
		seed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return a, fmt.Errorf("parse seed %q: %w", s, err)
		}
		a.seeds[j] = seed
	}
	i++

	// Parse maps
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse map header (e.g., "seed-to-soil map:")
		if !strings.HasSuffix(line, "map:") {
			return a, fmt.Errorf("parse map: expected map header, got %q", line)
		}
		mapName := strings.TrimSuffix(line, " map:")

		// Parse mappings until next empty line or EOF
		m := AlmanacMap{name: mapName}
		for i < len(lines) {
			line = strings.TrimSpace(lines[i])
			if line == "" {
				i++
				break
			}

			parts := strings.Fields(line)
			if len(parts) != 3 {
				return a, fmt.Errorf("parse mapping in %s: expected 3 numbers, got %d in %q", mapName, len(parts), line)
			}

			dst, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return a, fmt.Errorf("parse mapping in %s, destination: %w", mapName, err)
			}
			src, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return a, fmt.Errorf("parse mapping in %s, source: %w", mapName, err)
			}
			len, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				return a, fmt.Errorf("parse mapping in %s, length: %w", mapName, err)
			}

			m.mappings = append(m.mappings, Mapping{
				dstStart: dst,
				srcStart: src,
				length:   len,
			})
			i++
		}

		if len(m.mappings) == 0 {
			return a, fmt.Errorf("parse map %s: no mappings found", mapName)
		}
		a.maps = append(a.maps, m)
	}

	if len(a.maps) == 0 {
		return a, fmt.Errorf("parse almanac: no maps found")
	}

	return a, nil
}

// solve returns (part1, part2, error).
func solve(input string) (int64, int64, error) {
	almanac, err := ParseAlmanac(input)
	if err != nil {
		return 0, 0, err
	}

	// Part 1: Find minimum location for individual seeds
	var part1 int64 = -1
	for _, seed := range almanac.seeds {
		loc := almanac.LocationOf(seed)
		if part1 == -1 || loc < part1 {
			part1 = loc
		}
	}
	if part1 == -1 {
		return 0, 0, fmt.Errorf("solve part1: no seeds")
	}

	// Part 2: Treat seeds as pairs (start, length) and find minimum location
	if len(almanac.seeds)%2 != 0 {
		return 0, 0, fmt.Errorf("solve part2: seed count must be even, got %d", len(almanac.seeds))
	}

	seedRanges := make([]Range, 0, len(almanac.seeds)/2)
	for i := 0; i < len(almanac.seeds); i += 2 {
		seedRanges = append(seedRanges, Range{
			start:  almanac.seeds[i],
			length: almanac.seeds[i+1],
		})
	}

	part2 := almanac.LocationOfRanges(seedRanges)

	return part1, part2, nil
}
