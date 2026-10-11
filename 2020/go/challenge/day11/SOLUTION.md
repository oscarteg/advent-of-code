# Day 11: Seating System

Simulate people choosing seats in a waiting area.

## Part A
Apply seating rules based on **adjacent seats** (8 neighbors):
- Empty seat (L) becomes occupied (#) if no adjacent seats are occupied
- Occupied seat (#) becomes empty (L) if 4+ adjacent seats are occupied
- Floor (.) never changes

Repeat until stable, then count occupied seats.

## Part B
Apply seating rules based on **line-of-sight** visibility:
- Empty seat becomes occupied if no occupied seats are visible in any direction
- Occupied seat becomes empty if 5+ occupied seats are visible

Repeat until stable, then count occupied seats.

## Complexity
- **Time:** O(rows × cols × iterations × 8) per part
- **Space:** O(rows × cols) for grid storage

## Run & Test
```bash
cd 2020/go
go test -race ./challenge/day11
go run . 11 < ../rust/input/2020/day11.txt
```
