# Day 5: If You Give A Seed A Fertilizer

## Approach

**Part 1:** For each seed, trace it through the transformation pipeline sequentially. A value maps if it falls within a source range; otherwise it passes through unchanged. Return the minimum final location.

**Part 2:** Treat seed pairs as ranges [start, start+length). Rather than enumerate billions of seeds, work with intervals. For each map stage, split incoming ranges at mapping boundaries, transform overlaps, and carry the resulting ranges into the next stage. Return the minimum start value of the final ranges.

## Algorithm

The key insight for Part 2 is interval arithmetic:
- For each range [current, end) and each mapping, find overlaps and transform them
- Unmapped gaps between mappings stay unchanged
- After each map stage, all resulting ranges are carried forward

This avoids enumerating individual seeds and runs efficiently even for huge seed counts.

## Complexity

- **Part 1:** O(S × M) where S = number of seeds, M = total mapping entries
- **Part 2:** O(R × M) where R = number of range fragments (grows with boundary crossings, not seed count)

## Commands

```sh
cd 2023/go
go test ./day05
go run ./day05 < day05/testdata/example.txt  # 35, 46
```
