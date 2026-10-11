# Day 7: No Space Left On Device

## Part 1
Replay the `cd`/`ls` transcript with a path stack. Each file size is added to
every ancestor directory currently on the stack (keyed by absolute path), so
every directory's total size (including subdirectories) is known after a
single pass. Sum the sizes of directories at most 100,000.

Complexity: O(n * d) where n is the number of lines and d is max directory
depth (each file touches every ancestor); O(n) directories stored.

## Part 2
Compute used space as the root directory's size, derive unused space from the
70,000,000-byte disk, and the amount still needed to reach 30,000,000 free.
Pick the smallest directory whose size is at least that amount.

Complexity: O(n) to scan directory sizes for the minimum.

## Run
```bash
cargo run -q --bin day7 < input/day7-example.txt
```

## Test
```bash
cargo test --bin day7
```
