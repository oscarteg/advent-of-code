package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidTranscript = errors.New("invalid terminal transcript")

const (
	diskSpace           uint64 = 70_000_000
	requiredUnused      uint64 = 30_000_000
	smallDirectoryLimit uint64 = 100_000
)

func pathKey(stack []string) string {
	if len(stack) == 0 {
		return "/"
	}
	return "/" + strings.Join(stack, "/")
}

func directorySizes(input string) (map[string]uint64, error) {
	sizes := make(map[string]uint64)
	var stack []string

	for _, raw := range strings.Split(input, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || line == "$ ls" {
			continue
		}

		if target, ok := strings.CutPrefix(line, "$ cd "); ok {
			switch target {
			case "/":
				stack = stack[:0]
			case "..":
				if len(stack) == 0 {
					return nil, fmt.Errorf("%w: cd .. attempted above root", errInvalidTranscript)
				}
				stack = stack[:len(stack)-1]
			default:
				stack = append(stack, target)
			}
			sizes[pathKey(stack)] += 0
			continue
		}

		if name, ok := strings.CutPrefix(line, "dir "); ok {
			if name == "" {
				return nil, fmt.Errorf("%w: directory entry missing a name", errInvalidTranscript)
			}
			continue
		}

		sizeStr, name, found := strings.Cut(line, " ")
		if !found || name == "" {
			return nil, fmt.Errorf("%w: malformed listing line %q", errInvalidTranscript, line)
		}
		size, err := strconv.ParseUint(sizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid file size %q: %w", errInvalidTranscript, sizeStr, err)
		}

		for depth := 0; depth <= len(stack); depth++ {
			sizes[pathKey(stack[:depth])] += size
		}
	}

	if len(sizes) == 0 {
		return nil, fmt.Errorf("%w: transcript contains no directories", errInvalidTranscript)
	}
	return sizes, nil
}

func solveWith(input string, disk, unusedNeeded, smallLimit uint64) (uint64, uint64, error) {
	sizes, err := directorySizes(input)
	if err != nil {
		return 0, 0, err
	}

	var part1 uint64
	for _, size := range sizes {
		if size <= smallLimit {
			part1 += size
		}
	}

	used, ok := sizes["/"]
	if !ok {
		return 0, 0, fmt.Errorf("%w: missing root directory", errInvalidTranscript)
	}
	if used > disk {
		return 0, 0, fmt.Errorf("%w: used space exceeds total disk space", errInvalidTranscript)
	}
	unused := disk - used
	var needToFree uint64
	if unusedNeeded > unused {
		needToFree = unusedNeeded - unused
	}

	part2, found := uint64(0), false
	for _, size := range sizes {
		if size >= needToFree && (!found || size < part2) {
			part2 = size
			found = true
		}
	}
	if !found {
		return 0, 0, fmt.Errorf("%w: no directory is large enough to delete", errInvalidTranscript)
	}

	return part1, part2, nil
}

func solve(input string) (uint64, uint64, error) {
	return solveWith(input, diskSpace, requiredUnused, smallDirectoryLimit)
}
