# Day 3: Rucksack Reorganization

## Part 1
Each rucksack's contents are split into two equal halves (compartments). Find the item type
common to both halves and sum its priority (a-z = 1-26, A-Z = 27-52) across all rucksacks.

**Approach:** Split each line in half, intersect the two halves as sets, sum priorities.
**Complexity:** O(n) where n is the total number of items.

## Part 2
Rucksacks are grouped in threes; each group of three elves shares exactly one common badge
item. Sum the badge priorities across all groups.

**Approach:** Chunk lines by 3, intersect the three lines as sets, sum priorities.
**Complexity:** O(n) where n is the total number of items.

## Run
```bash
cargo run -q --bin day3 < input/day3-example.txt
```

## Test
```bash
cargo test --bin day3
```
