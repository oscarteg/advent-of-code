package main

import (
	"errors"
	"fmt"
)

var errInvalidDiskMap = errors.New("invalid disk map")

const freeID = int64(-1)

type segment struct {
	id     int64
	start  int64
	length int64
}

func parseSegments(input string) ([]segment, error) {
	digits := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '\n' || c == '\r' || c == ' ' {
			continue
		}
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("%w: invalid character %q", errInvalidDiskMap, c)
		}
		digits = append(digits, c)
	}
	if len(digits) == 0 {
		return nil, fmt.Errorf("%w: empty disk map", errInvalidDiskMap)
	}

	segments := make([]segment, 0, len(digits))
	var position, fileID int64
	for i, digit := range digits {
		length := int64(digit - '0')
		isFile := i%2 == 0
		if length > 0 {
			id := freeID
			if isFile {
				id = fileID
			}
			segments = append(segments, segment{id: id, start: position, length: length})
			position += length
		}
		if isFile {
			fileID++
		}
	}
	return segments, nil
}

func totalLength(segments []segment) int64 {
	var total int64
	for _, s := range segments {
		total += s.length
	}
	return total
}

// solvePart1 fragments files block by block, so it is simulated on an explicit block array.
func solvePart1(segments []segment) int64 {
	blocks := make([]int64, totalLength(segments))
	for i := range blocks {
		blocks[i] = freeID
	}
	for _, s := range segments {
		for i := int64(0); i < s.length; i++ {
			blocks[s.start+i] = s.id
		}
	}

	left, right := 0, len(blocks)
	for {
		for left < len(blocks) && blocks[left] != freeID {
			left++
		}
		for right > 0 && blocks[right-1] == freeID {
			right--
		}
		if left >= right {
			break
		}
		blocks[left] = blocks[right-1]
		blocks[right-1] = freeID
		left++
		right--
	}

	var checksum int64
	for i, block := range blocks {
		if block != freeID {
			checksum += int64(i) * block
		}
	}
	return checksum
}

// solvePart2 moves whole files, highest id first, into the leftmost free span that
// fits. Free spans stay sorted by start as they only shrink from the left, so a
// single left-to-right scan per file finds the leftmost fit.
func solvePart2(segments []segment) int64 {
	var files, freeSpans []segment
	for _, s := range segments {
		if s.id == freeID {
			freeSpans = append(freeSpans, s)
		} else {
			files = append(files, s)
		}
	}

	var checksum int64
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		position := file.start
		for j := range freeSpans {
			free := &freeSpans[j]
			if free.start >= file.start {
				break
			}
			if free.length >= file.length {
				position = free.start
				free.start += file.length
				free.length -= file.length
				break
			}
		}
		for k := int64(0); k < file.length; k++ {
			checksum += (position + k) * file.id
		}
	}
	return checksum
}

func solve(input string) (part1, part2 int64, err error) {
	segments, err := parseSegments(input)
	if err != nil {
		return 0, 0, err
	}
	return solvePart1(segments), solvePart2(segments), nil
}
