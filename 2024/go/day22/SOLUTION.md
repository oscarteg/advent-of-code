# Day 22: Monkey Market

Evolve each buyer's secret 2000 times via a fixed mix-and-prune sequence
(multiply/divide, xor, mod 2^24); part 1 sums the 2000th secret per buyer.

## Approach

Part 2: the secret's last digit is the banana price. Track the last four
price deltas as a key and, per buyer, record the price only the first time
each key appears. Summing each key's price across all buyers and taking the
max gives the best banana total over every 4-delta selling sequence.

## Complexity

O(buyers x 2000) time and space.

## Run / test

```bash
cd 2024/go && go test -race -count=1 ./day22
go run ./day22 < day22/testdata/example.txt
```
