# Day 13: Shuttle Search

## Part 1
For each bus id, compute the wait until the next departure at or after the
earliest timestamp as `(-earliest).rem_euclid(id)`, then pick the bus with
the minimum wait and return `id * wait`. O(n) in the number of buses.

## Part 2
Find the earliest timestamp where each bus departs `offset` minutes after it,
using an incremental Chinese remainder theorem: fold buses in one at a time,
advancing the running timestamp by the combined period (product of ids seen
so far) until it satisfies the new bus's constraint too. O(n) bus-steps per
bus since AoC bus ids are pairwise coprime (prime), independent of timestamp
magnitude (no brute-force search over timestamps).

## Run and test
```bash
cargo test --bin day13
cargo run -q --bin day13 < input/2020/day13-example.txt
```
