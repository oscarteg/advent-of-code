# Day 9: Disk Fragmenter (Zig)

A dense disk map alternates file-length and free-length digits, starting
with the first file (id 0). Compact the disk and compute a checksum: the sum
of `position * file_id` over every used block.

## Approach

Expand the map into typed segments `{id, start, length}` (`id = -1` for
free space). Part 1 fragments files block by block: walk a left pointer to
the next free block and a right pointer to the last used block, swapping
until they meet, then checksum the resulting block array.

Part 2 moves whole files without fragmenting them, highest file id first,
into the leftmost free span that fits. Free spans only shrink from the left
as files land in them, so they stay sorted by start and a single scan per
file finds the leftmost fit (or leaves the file in place).

## Complexity

Part 1: O(n) in the number of blocks. Part 2: O(files x free spans).
Space: O(n).

## Run / test

```bash
cd 2024/zig
zig test day9.zig
zig run day9.zig < examples/day9.txt
```
