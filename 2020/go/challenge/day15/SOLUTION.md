# Day 15: Rambunctious Recitation

## Part A & B

Simulate the memory game: each turn, say the number of turns since the last
time the previous number was spoken, or 0 if it's new. Both parts run the same
`play` simulation, just to a different turn count (2020 vs 30,000,000).

`play` uses a flat `[]int` indexed by spoken number to store the last turn each
number was said (0 = never seen), keeping each turn O(1). Time complexity is
O(turns), space is O(max spoken value).

## Run & test

```bash
go run . 15 < ../rust/input/2020/day15.txt
go test -race ./challenge/day15
```
