# Day 5: If You Give A Seed A Fertilizer

- **Part 1:** Follow each seed through the maps. A matching source interval
  translates the number; unmatched numbers stay unchanged. Take the minimum.
- **Part 2:** Treat seed pairs as half-open ranges. Split at mapping boundaries,
  translate overlaps, and retain unmapped gaps. Carry the resulting ranges into
  the next map. The smallest final range start is the answer.
- Runtime depends on mapping rows and range fragments, not the number of seeds.
  Part 1 uses `O(seeds * total rows)` time; part 2 scans each stage's rows for
  every incoming fragment and stores the outgoing fragments.

From `2023/rust`:

```sh
cargo run -q --bin day05 < examples/day05.txt # 35
cargo run -q --bin day05 < examples/day05.txt # 46
cargo test day05::
```

Use the same commands with your puzzle input redirected into stdin.
