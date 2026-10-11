# Day 10: Cathode-Ray Tube

## Approach

**Part 1:** Simulate a CPU with register X (initially 1). Process instructions:
- `noop`: 1 cycle, no change
- `addx V`: 2 cycles, add V to X after both cycles complete

Calculate signal strength (cycle × X) at cycles 20, 60, 100, 140, 180, 220 and sum them.

**Part 2:** Render a 6×40 CRT display:
- Sprite is 3 pixels wide centered at X (covers X-1, X, X+1)
- During each cycle, a pixel lights up if the current column is covered by the sprite
- Output the 6 rows as a multi-line string

## Complexity

- Time: O(n) where n is the number of instructions
- Space: O(240) for storing all cycle X values

## Commands

Run:
```bash
cd 2022/zig
zig run src/day_10.zig < src/examples/day_10.txt
```

Test:
```bash
cd 2022/zig
zig test src/day_10.zig
```

Verify:
```bash
cd 2022/zig
zig fmt --check src/day_10.zig && zig test src/day_10.zig && zig run src/day_10.zig < src/examples/day_10.txt
```
