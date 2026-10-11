# Day 11: Reactor

## Problem

Every input line `device: outputs...` names a device and the devices it
feeds, together forming a directed graph.

**Part 1:** count the distinct paths from `you` to `out`.

**Part 2:** count the distinct paths from `svr` to `out` that pass through
both `dac` and `fft`.

---

*Solution spoilers below.*

## Approach

- **Parsing:** every `device: outputs...` line becomes a map from a device to the devices it feeds. A line without `:` or without a device name is an error, and so is an input without any device.
- **Part 1:** a memoized depth-first search counts the paths from `you` to `out`. Reaching the target counts as one path; any other device's count is the sum of its outputs' counts.
- **Part 2:** a path through both `dac` and `fft` meets them in one of two orders. For each order, `pathsVia` multiplies the counts of its three segments (`svr` to the first waypoint, first to second, second to `out`), each with its own memo. The answer is the sum of both orders. The graph is acyclic, so at most one order has any paths and no path is counted twice.

## Complexity

O(V + E) time and O(V) memo space per segment; seven segments in total, so O(V + E) overall.

## Commands

```bash
cd 2025
zig test src/day11.zig
zig run src/day11.zig < src/examples/day11.txt
```

```bash
cd 2025/go
go test -race -count=1 ./day11
go run ./day11 < day11/testdata/example.txt
```

The part 1 example has no `svr`, so both runners print `Part 1: 5` and `Part 2: 0` for it. The official part 2 example, where `fft` comes before `dac`, gives 2 and is covered by the tests in both languages.
