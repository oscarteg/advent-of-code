# Day 4: Passport Processing

## Part 1
Split the input on blank lines into passport blocks, then each block's whitespace-separated
`key:value` pairs into a `HashMap<String, String>`. Count passports that contain all seven
required fields (`cid` is optional). Complexity: O(n) in input size.

## Part 2
Same parsing, but each required field is additionally validated against its format/range rule
(year bounds, `cm`/`in` height ranges, hex hair color, enumerated eye colors, 9-digit passport
ID) using two small regexes plus numeric range checks. Complexity: O(n).

## Run
```bash
cargo run -q --bin day4 < input/2020/day4-example.txt
```

## Test
```bash
cargo test --bin day4
```
