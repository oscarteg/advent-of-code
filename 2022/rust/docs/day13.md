# Day 13: Distress Signal

## Approach

**Part 1:** Compare pairs of packets in order, using a custom comparison:
- Integers compare numerically
- Lists compare element-by-element; shorter list is "less"
- When comparing int vs list, wrap the int and compare

Sum the 1-indexed pair numbers that are in correct order.

**Part 2:** Parse all packets (skip blank lines), add dividers `[[2]]` and `[[6]]`, sort all using the same comparison, then find the product of the 1-indexed divider positions.

## Complexity

- **Time:** O(n log n) for sorting in part 2; O(n * d) for comparison where d is max depth
- **Space:** O(n) for packet storage and sorting

## Run & Test

```bash
cd 2022/rust
rustfmt --edition 2021 --check src/days/day13.rs src/bin/day13.rs
cargo test --bin day13
cargo clippy --bin day13 -- -D warnings
cargo run -q --bin day13 < input/day13-example.txt
```

Expected output: `Part 1: 13` and `Part 2: 140`
