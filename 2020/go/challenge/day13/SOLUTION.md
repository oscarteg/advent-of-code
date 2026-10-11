# Day 13: Shuttle Search

## Part A
For each bus id, wait = `(id - earliest % id) % id`; pick the bus with the
smallest wait and return `id * wait`. O(n) buses.

## Part B
Incremental Chinese remainder theorem: fold buses in one at a time, advancing
a running timestamp by the combined period (product of ids seen so far)
until `(timestamp + offset) % id == 0`. O(n) bus-steps, no brute-force search
over timestamps.

## Run and test
```bash
go run . 13 < ../rust/input/2020/day13.txt
go test -race -count=1 ./challenge/day13
```
