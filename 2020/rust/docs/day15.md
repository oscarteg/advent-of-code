# Day 15: Rambunctious Recitation

## Part 1 & 2

Play the memory game: each turn, say the number of turns since the last number
spoken was previously said, or 0 if it's the first time. Both parts are the same
simulation, just run to a different turn count (2020 vs 30,000,000).

A flat `Vec<u32>` indexed by spoken number stores the last turn each number was
said (0 = never seen), making each turn O(1) instead of scanning history. Time
complexity is O(turns), space is O(max spoken value).

## Run & test

```bash
cargo test --bin day15
cargo run -q --bin day15 < input/2020/day15-example.txt
```
