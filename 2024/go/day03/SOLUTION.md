# Day 3: Mull It Over

## Approach

**Part 1:** Find all valid `mul(X,Y)` patterns in the input, where X and Y are 1-3 digit numbers. Sum all the products.

**Part 2:** Same as part 1, but additionally track `do()` and `don't()` instructions. The execution state starts enabled. Only sum products when in the enabled state.

## Implementation

Scan the input string sequentially, finding the next occurrence of `mul(`, `do()`, or `don't()`. Process them in order:
- For `mul(X,Y)`: Extract X and Y manually by parsing digits, compute the product, add to part 1 total. If enabled, add to part 2 total.
- For `do()`: Set enabled = true.
- For `don't()`: Set enabled = false.

## Complexity

- Time: O(n) where n is the length of the input (linear scan).
- Space: O(1) additional space beyond input storage.

## Commands

```bash
go test -race ./day03
go run ./day03 < day03/testdata/example.txt    # part 1 example: 161
go run ./day03 < day03/testdata/example2.txt   # part 2 example: 48
```
