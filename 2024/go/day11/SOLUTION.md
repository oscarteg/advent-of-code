# Day 11: Plutonian Pebbles

Each stone transforms per the blink rule: `0` becomes `1`; an even digit
count splits into two stones (the halves, no leading zeros kept); otherwise
the value is multiplied by 2024. Track a count per distinct stone value
instead of the list itself, since the sequence grows exponentially but the
set of distinct values stays small. Part 1 reads the total at 25 blinks,
part 2 at 75, from the same run.

## Complexity

O(blinks x distinct values) time and space.

## Run / test

```bash
cd 2024/go && go test -race -count=1 ./day11
go run ./day11 < day11/testdata/example.txt
```
