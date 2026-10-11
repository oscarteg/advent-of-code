# Day 13: Transparent Origami

## Part 1
Parse dots into a `HashSet<(i32, i32)>` and fold instructions into `Fold::X`/`Fold::Y`.
Apply only the first fold: every dot past the fold line reflects across it
(`2 * line - coord`), dots are deduplicated by the set. Answer is the resulting
dot count. Complexity: O(n) in the number of dots.

## Part 2
Apply every remaining fold in order, then render the final dot set as a grid of
`#`/`.` lines sized to the surviving bounding box. Complexity: O(n * f) for
n dots and f folds.

## Run
```bash
cargo run -q --bin day13 < input/day13-example.txt
```

## Test
```bash
cargo test --bin day13
```
