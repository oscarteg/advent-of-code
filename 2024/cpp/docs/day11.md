# Day 11: Plutonian Pebbles (C++)

Each stone transforms per the blink rule: `0` becomes `1`; an even number of
digits splits into two stones (left/right halves, no leading zeros kept);
otherwise the value is multiplied by 2024. The sequence explodes
combinatorially, so track a count per distinct stone value instead of the
list of stones, and merge counts for values that recur. Part 1 blinks 25
times, part 2 blinks 75 times; both read off the total count at their
checkpoint from the same run.

## Complexity

O(blinks x distinct values) time and space; the number of distinct values
stays small (a few thousand) because the rule only maps into a bounded set
of derived numbers.

## Run / test

```bash
cd 2024/cpp
make test day=day11
make -s run day=day11 < examples/day11.txt
```
