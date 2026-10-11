# Day 10: Cathode-Ray Tube

## Approach

**Part 1:** Simulate a CPU with a single register X (initially 1). Process two instruction types:
- `noop`: 1 cycle, no operation
- `addx V`: 2 cycles, add V to X

Accumulate signal strength (cycle × X) at specific cycles (20, 60, 100, 140, 180, 220) and sum them.

**Part 2:** Simulate a CRT display 40 pixels wide and 6 pixels tall. During each cycle, draw a pixel:
- The sprite is 3 pixels wide centered at X (covers X-1, X, X+1)
- A pixel lights up if the current drawing position overlaps the sprite
- Output as a 6×40 grid using '#' (lit) and '.' (dark)

## Complexity

- Time: O(n) where n is the number of instructions
- Space: O(1) for tracking CPU state (O(240) for storing all x_values)

## Commands

Run:
```bash
cd 2022/rust
cargo run -q --bin day10 < input/day10-example.txt
```

Test:
```bash
cd 2022/rust
cargo test --bin day10
```

Verify:
```bash
cd 2022/rust
rustfmt --edition 2021 --check src/days/day10.rs src/bin/day10.rs && cargo test --bin day10 && cargo clippy --bin day10 -- -D warnings && cargo run -q --bin day10 < input/day10-example.txt
```
