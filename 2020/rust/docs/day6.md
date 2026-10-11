# Day 6: Custom Customs

## Part 1
Split the input into groups on blank lines. For each group, collect the set
union of letters answered by any person and sum the set sizes across groups.

## Part 2
Same grouping, but take the set intersection of letters answered by every
person in the group instead of the union, then sum the sizes.

## Complexity
O(n) in the total input length: each line is parsed into a `HashSet<char>`
(at most 26 entries) once, and groups are folded with union/intersection.

## Run and test
```bash
cargo test --bin day6
cargo run -q --bin day6 < input/2020/day6-example.txt
```
