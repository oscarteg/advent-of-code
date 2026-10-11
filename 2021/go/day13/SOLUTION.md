# Day 13: Transparent Origami

## Part 1
Parse dots into a `map[point]struct{}` and fold instructions into `fold{axis, value}`.
Apply only the first fold: every dot past the fold line reflects across it
(`2*value - coord`), dots are deduplicated by the map. Answer is the resulting
dot count. Complexity: O(n) in the number of dots.

## Part 2
Apply every remaining fold in order, then render the final dot set as a grid of
`#`/`.` lines sized to the surviving bounding box. Complexity: O(n * f) for
n dots and f folds.

## Run
```bash
go run ./day13 < day13/testdata/example.txt
```

## Test
```bash
go test -race -count=1 ./day13
```
