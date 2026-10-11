# Day 20: Grove Positioning System

## Part 1

Parse numbers, then mix the circular list once: for each number in original
order, remove it from an `order` slice of original indices and reinsert it
`value` slots forward/backward, wrapping modulo `len - 1` (removing an
element temporarily shortens the cycle it moves within). Indices, not
values, are tracked because values can repeat. Sum the values 1000, 2000,
and 3000 positions after 0, wrapping around.

Complexity: O(n^2) per mix round (linear index lookup + slice splice), O(n)
rounds for part 2.

## Part 2

Multiply every number by the decryption key (811589153), mix 10 times
instead of 1, reuse the same summing logic.

## Run

```bash
go run ./day20 < day20/testdata/example.txt
```

## Test

```bash
go test -race -count=1 ./day20
```
