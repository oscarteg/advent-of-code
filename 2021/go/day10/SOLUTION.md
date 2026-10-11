# Day 10: Syntax Scoring

## Approach

**Part 1:** Parse navigation subsystem lines character by character, tracking open brackets on a stack. When encountering a closing bracket, check if it matches the most recently opened bracket. If not, record the mismatch score and move to the next line. Sum all error scores.

**Part 2:** For incomplete lines (those without syntax errors), simulate closing all open brackets in reverse order, computing an autocomplete score for each. Return the median score.

## Complexity

- **Time:** O(n·m) where n is the number of lines and m is average line length
- **Space:** O(m) for the bracket stack

## Commands

```bash
cd 2021/go
go test ./day10
go run ./day10 < day10/testdata/example.txt
```

Expected output:
```
Part 1: 26397
Part 2: 288957
```
