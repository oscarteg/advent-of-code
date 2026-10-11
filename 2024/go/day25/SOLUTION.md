# Day 25: Code Chronicle

Single part: count lock/key pairs whose pins cannot overlap.

## Idea

Parse each 7x5 schematic: a lock has a full top row, a key a full bottom
row. Pin height per column is the `#` count minus the full row. A lock and
key fit if every column's combined height stays within the 5-row gap
between the two full rows. Count all fitting pairs by brute force.

## Complexity

O(n) parsing; O(L*K) fit checks for L locks and K keys.

## Run / test

```
cd 2024/go
go test ./day25
go run ./day25 < day25/testdata/example.txt
```
