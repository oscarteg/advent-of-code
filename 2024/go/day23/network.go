package main

import (
	"fmt"
	"sort"
	"strings"
)

func solve(input string) (int, string, error) {
	graph := make(map[string]map[string]bool)

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return 0, "", fmt.Errorf("invalid edge format: %s", line)
		}
		a, b := parts[0], parts[1]

		if graph[a] == nil {
			graph[a] = make(map[string]bool)
		}
		if graph[b] == nil {
			graph[b] = make(map[string]bool)
		}
		graph[a][b] = true
		graph[b][a] = true
	}

	triangles := make(map[[3]string]bool)
	for node := range graph {
		neighbors := make([]string, 0)
		for neighbor := range graph[node] {
			neighbors = append(neighbors, neighbor)
		}
		sort.Strings(neighbors)

		for i := 0; i < len(neighbors); i++ {
			for j := i + 1; j < len(neighbors); j++ {
				n1, n2 := neighbors[i], neighbors[j]
				if graph[n1][n2] {
					tri := [3]string{node, n1, n2}
					sort.Strings(tri[:])
					triangles[tri] = true
				}
			}
		}
	}

	part1 := 0
	for tri := range triangles {
		for _, node := range tri {
			if len(node) > 0 && node[0] == 't' {
				part1++
				break
			}
		}
	}

	nodes := make([]string, 0)
	for node := range graph {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	var maxClique []string
	n := len(nodes)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				for l := k + 1; l < n; l++ {
					cand := []string{nodes[i], nodes[j], nodes[k], nodes[l]}
					if isClique(cand, graph) && len(cand) > len(maxClique) {
						maxClique = make([]string, len(cand))
						copy(maxClique, cand)
					}
				}
			}
		}
	}

	if len(maxClique) == 0 {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				for k := j + 1; k < n; k++ {
					cand := []string{nodes[i], nodes[j], nodes[k]}
					if isClique(cand, graph) && len(cand) > len(maxClique) {
						maxClique = make([]string, len(cand))
						copy(maxClique, cand)
					}
				}
			}
		}
	}

	sort.Strings(maxClique)
	part2 := strings.Join(maxClique, ",")

	return part1, part2, nil
}

func isClique(nodes []string, graph map[string]map[string]bool) bool {
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if !graph[nodes[i]][nodes[j]] {
				return false
			}
		}
	}
	return true
}
