# Day 14: Extended Polymerization

## Problem

Start with a polymer template and apply pair insertion rules repeatedly. Each rule specifies what character to insert between two characters. After N steps, count the most and least common characters.

## Solution

**Part 1 (10 steps) and Part 2 (40 steps):**
- Track character frequencies as we apply rules
- Instead of building the full string (which would be massive), count pair occurrences
- When a pair `AB` has rule `AB -> C`, it becomes `ACB`, creating pairs `AC` and `CB`
- Update character counts incrementally: each insertion adds 1 to the insertion character
- Answer: `max_count - min_count`

## Complexity

- Time: O(R * S) where R is the number of rules and S is the number of steps
- Space: O(R) for storing rules and pair counts

## Commands

```bash
cargo test --bin day14
cargo run -q --bin day14 < input/day14-example.txt
```

Expected output:
```
Part 1: 1588
Part 2: 2188189693529
```
