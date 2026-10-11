# Day 8: Resonant Collinearity (C++)

## Problem

Given a grid of antennas (each marked with a frequency: lowercase letter, uppercase letter, or digit), find antinodes:
- An antinode occurs on the line between two antennas of the same frequency
- Specifically, for antennas at positions p1 and p2, antinodes are at p1 + (p1 - p2) and p2 + (p2 - p1)
- Count only antinodes within map bounds

Part 1: Count unique antinode positions (exactly two per antenna pair)
Part 2: Count all collinear points on the line through each pair of same-frequency antennas

## Approach

Part 1:
- Group antennas by frequency
- For each pair of same-frequency antennas, compute the two antinode positions
- Keep only those within bounds and count unique positions

Part 2:
- For each pair of same-frequency antennas, compute the direction (unit step after GCD normalization)
- Walk the line in both directions from each antenna, adding all points within bounds
- Count unique positions (includes both antenna positions and all collinear antinodes)

## Complexity

- Time: O(n² × m) where n is number of antennas per frequency, m is grid width/height for line walking
- Space: O(m²) for the antinode set

## Example Answers

Part 1: 14
Part 2: 34

## Run / test

```bash
cd 2024/cpp
make test day=day8
make -s run day=day8 < examples/day8.txt
```
