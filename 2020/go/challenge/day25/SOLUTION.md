# Day 25: Combo Breaker

## Approach

Single part. Parse the card and door public keys. Brute-force the card's loop
size by repeatedly multiplying a running value by the subject number 7,
modulo 20201227, until it matches the card's public key. Apply that loop size
to the door's public key (the same transform) to get the shared encryption
key.

## Complexity

- Time: O(L) where L is the loop size found (bounded by the modulus).
- Space: O(1).

## Commands

```bash
go test -race -count=1 ./challenge/day25
go run . 25 < ../rust/input/2020/day25.txt
```
