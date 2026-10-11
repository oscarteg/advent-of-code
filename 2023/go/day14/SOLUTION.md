# Day 14: Parabolic Reflector Dish

## Part 1: Tilt North

The dish is a grid with rolling rocks ('O'), fixed rocks ('#'), and empty space ('.'). When tilting north, rolling rocks move as far up as possible until hitting a fixed rock or the edge.

The load is calculated as the sum of (number of rocks in row) × (distance from bottom).

Algorithm: For each column, scan from top, accumulating rocks and placing them as high as possible.

Time: O(rows × cols) | Space: O(1) excluding grid.

## Part 2: Billion Cycles

Run spin cycles (north, west, south, east) 1,000,000,000 times. Since the grid state repeats, detect the cycle and calculate the final state without simulating all iterations.

Algorithm: Store grid states in a map. When we see a repeated state, calculate how many remaining cycles are needed modulo the cycle length.

Time: O(cycle_length × rows × cols) | Space: O(cycle_length × rows × cols) for the map.

## Commands

```bash
cd 2023/go
go test ./day14
go run ./day14 < day14/testdata/example.txt
```
