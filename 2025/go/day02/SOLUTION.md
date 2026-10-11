# Day 2: Gift Shop

Input: one line of comma-separated inclusive ID ranges `a-b`.

## Part 1

An ID counts if its decimal digits are a shorter sequence repeated exactly
twice (`55`, `6464`, `123123`): even digit count, first half equals second
half. Sum the counted IDs across every range.

## Part 2

An ID counts if its digits are a shorter sequence repeated two or more
times (`123123123`, `1111111` also count). For every candidate period that
evenly divides the digit count, check whether that period tiles the whole
string.

## Complexity

`O(N * D^2)` where `N` is the number of IDs spanned by all ranges and `D`
is the digit count of the largest ID: each ID is checked against up to `D`
candidate periods, each an `O(D)` scan.

## Run / test

```bash
cd 2025/go
go test ./day02
go run ./day02 < day02/testdata/example.txt
```
