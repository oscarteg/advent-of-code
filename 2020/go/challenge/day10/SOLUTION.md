# Day 10: Adapter Array

## Part A
Sort the adapter joltages, prepend the outlet (0) and append the device
(highest adapter + 3). Walk the sorted chain counting the 1-jolt and 3-jolt
gaps between consecutive values, then multiply the counts.

Complexity: O(n log n) for the sort, O(n) for the scan.

## Part B
Count distinct valid chains with dynamic programming over joltage values:
`ways[j]` is the number of ways to reach joltage `j`, equal to the sum of
`ways[j-1]`, `ways[j-2]`, `ways[j-3]` for whichever of those joltages exist
in the chain. `ways[0] = 1` (the outlet), and the answer is `ways[device]`.

Complexity: O(n) time and space, where n is the device's joltage.

## Run and test
```bash
go run . 10 < ../rust/input/2020/day10.txt
go test -race ./challenge/day10
```
