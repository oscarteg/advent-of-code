package main

import (
	"fmt"
	"strings"
)

func solve(input string) (int, int, error) {
	graph, err := parseGraph(input)
	if err != nil {
		return 0, 0, err
	}

	part1 := countPaths(graph, false)
	part2 := countPaths(graph, true)

	return part1, part2, nil
}

func parseGraph(input string) (map[string][]string, error) {
	graph := make(map[string][]string)

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected edge format 'a-b', got '%s'", line)
		}

		a, b := parts[0], parts[1]
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	return graph, nil
}

func isSmallCave(cave string) bool {
	for _, c := range cave {
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func countPaths(graph map[string][]string, allowDouble bool) int {
	visited := make(map[string]bool)
	return dfs(graph, "start", visited, allowDouble, allowDouble)
}

func dfs(
	graph map[string][]string,
	current string,
	visited map[string]bool,
	allowDouble bool,
	doubleAvailable bool,
) int {
	if current == "end" {
		return 1
	}

	if isSmallCave(current) {
		visited[current] = true
	}

	neighbors := graph[current]
	count := 0

	for _, neighbor := range neighbors {
		alreadyVisited := visited[neighbor]

		if !alreadyVisited {
			newVisited := copyVisited(visited)
			count += dfs(graph, neighbor, newVisited, allowDouble, doubleAvailable)
		} else if allowDouble && doubleAvailable && neighbor != "start" {
			newVisited := copyVisited(visited)
			count += dfs(graph, neighbor, newVisited, allowDouble, false)
		}
	}

	return count
}

func copyVisited(visited map[string]bool) map[string]bool {
	newVisited := make(map[string]bool)
	for k, v := range visited {
		newVisited[k] = v
	}
	return newVisited
}
