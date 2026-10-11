# Day 3: Mull It Over (C++)

## Approach

**Part 1:** Find all valid `mul(X,Y)` patterns in the input, where X and Y are 1-3 digit numbers. Sum all the products.

**Part 2:** Same as part 1, but additionally track `do()` and `don't()` instructions. The execution state starts enabled. Only sum products when in the enabled state.

## Implementation

One regular expression, `mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`, walks the corrupted memory left to right (C++ `std::cregex_iterator`; the Zig version scans the same three tokens by hand). `do()`/`don't()` flip an `enabled` flag that carries across lines; every `mul` adds to part 1 and, while enabled, to part 2.

## Complexity

- Time: O(n) where n is the length of the input (linear scan with regex matching).
- Space: O(1) additional space beyond input storage.

## Run / test

```bash
cd 2024/cpp
make test day=day3
make -s run day=day3 < examples/day3.txt
make -s run day=day3 < examples/day3-part2.txt
```
