# Day 6: Custom Customs

## Part A
Split input lines into groups on blank lines. Each person's answers become a
set of letters; the group's answer set is the union of all people's sets.
Sum the union sizes across groups.

## Part B
Same grouping, but the group's answer set is the intersection of all
people's sets instead of the union.

## Complexity
O(n) in the total input length: each line builds a set of at most 26
letters, and groups are folded with union/intersection.

## Run and test
```bash
go run . 6 < ../rust/input/2020/day6.txt
go test -race -count=1 ./challenge/day6
```
