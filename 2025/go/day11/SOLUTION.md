# Day 11: Reactor

## Approach

- **Parsing:** `device: outputs...` lines become a `graph` (`map[string][]string`). A line without `:` or without a device name is an error, and so is empty input.
- **Part 1:** `countPaths` is a memoized depth-first search from `you` to `out`. The target counts as one path; any other device's count is the sum of its outputs' counts.
- **Part 2:** a path through both `dac` and `fft` meets them in one of two orders. `pathsVia` multiplies the segment counts for one order (`svr` -> first -> second -> `out`), with a fresh memo per segment, and the answer adds both orders. The graph is acyclic, so at most one order has any paths and no path is counted twice.

## Complexity

O(V + E) time and O(V) memo space per segment.

## Commands

```bash
cd 2025/go
go test -race -count=1 ./day11
go run ./day11 < day11/testdata/example.txt
```

The part 1 example prints `Part 1: 5` and `Part 2: 0` because it has no `svr`. The tests also cover the official part 2 example, which gives 2.
