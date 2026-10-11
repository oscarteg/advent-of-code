# Day 6: Wait For It

- **Part 1:** Each column is a race `(time, distance)`. For a hold time `h`,
  travelled distance is `h * (time - h)`. Count every `h` in `0..=time` that
  beats the record, then multiply the counts across races.
- **Part 2:** Join the digits on each line (ignoring spaces) into one big
  race, then reuse the same counting function.
- Brute-force counting is `O(time)` per race; `time` stays small enough
  (even the joined part 2 race) that scanning every hold time is fast.

From `2023/go`:

```sh
go run ./day06 < day06/testdata/example.txt
go test -race ./day06
go vet ./day06
```

The example prints exactly:

```text
Part 1: 288
Part 2: 71503
```
