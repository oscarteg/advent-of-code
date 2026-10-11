# Day 4: Passport Processing

## Part A
Split the input on blank lines into passport blocks, parse each block's whitespace-separated
`key:value` fields into a map. Count passports that have all seven required fields (`cid` is
optional). Complexity: O(n) in input size.

## Part B
Same parsing, but each required field is additionally checked against its format/range rule:
year bounds, `cm`/`in` height ranges, hex hair color, enumerated eye colors, and a 9-digit
passport ID, using two regexes plus numeric range checks. Complexity: O(n).

## Run
```bash
go run . 4 < ../rust/input/2020/day4.txt
```

## Test
```bash
go test -race -count=1 ./challenge/day4
```
