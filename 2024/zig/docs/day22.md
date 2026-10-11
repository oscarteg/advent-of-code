# Day 22: Monkey Market (Zig)

Each buyer's secret number evolves through a fixed pseudo-random sequence:
multiply by 64, xor, mod 2^24; divide by 32, xor, mod 2^24; multiply by
2048, xor, mod 2^24. Part 1 sums the 2000th secret per buyer.

## Approach

Part 2: each secret's last digit is its banana price. Track the last four
consecutive price deltas as a key (base-19 encoding, offsetting each delta
by 9), and for each buyer record the price only at the first occurrence of
each key. Summing a key's price across all buyers and taking the maximum
gives the best banana total over every possible 4-delta selling sequence.

## Complexity

O(buyers x 2000) time and space per buyer for the seen-sequence set; O(19^4)
distinct keys overall.

## Run / test

```bash
cd 2024/zig
zig test day22.zig
zig run day22.zig < examples/day22.txt
```
