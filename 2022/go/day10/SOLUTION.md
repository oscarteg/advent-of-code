# Day 10: Cathode-Ray Tube

## Approach

**Part 1:** Simulate a CPU with register X (initially 1). Process two instruction types:
- `noop`: 1 cycle, no operation
- `addx V`: 2 cycles, add V to X

Accumulate signal strength (cycle × X) at specific cycles (20, 60, 100, 140, 180, 220) and sum them.

**Part 2:** Simulate a CRT display 40 pixels wide and 6 pixels tall. During each cycle, draw a pixel:
- The sprite is 3 pixels wide centered at X (covers X-1, X, X+1)
- A pixel lights up if the current drawing position overlaps the sprite
- Output as a 6×40 grid using '#' (lit) and '.' (dark)

## Complexity

- Time: O(n) where n is the number of instructions
- Space: O(240) for storing all X values

## Commands

Run:
```bash
cd 2022/go
go run ./day10 < day10/testdata/example.txt
```

Test:
```bash
cd 2022/go
go test ./day10
```

Verify:
```bash
cd 2022/go
test -z "$(gofmt -l day10)" && go vet ./day10 && go test -race -count=1 ./day10 && go run ./day10 < day10/testdata/example.txt
```
