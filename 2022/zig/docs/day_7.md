# Day 7: No Space Left On Device

## Part 1
Replay the `cd`/`ls` transcript with a path stack, adding each file's size
to every ancestor directory currently on the stack (keyed by absolute
path), so every directory's total size is known after a single pass. Sum
the sizes of directories at most 100,000.

Complexity: O(n * d), n lines and d max directory depth; O(n) directories.

## Part 2
Compute used space as the root directory's size, derive unused space from
the 70,000,000-byte disk, and the amount still needed to reach 30,000,000
free. Pick the smallest directory whose size is at least that amount.

Complexity: O(n) to scan directory sizes for the minimum.

## Run
```bash
zig run src/day_7.zig < src/examples/day_7.txt
```

## Test
```bash
zig test src/day_7.zig
```
