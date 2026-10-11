# Day 11: Cosmic Expansion

## Part 1
Find all galaxy positions and identify empty rows/columns (those with no galaxies). Calculate the Manhattan distance between every pair of galaxies, but count each empty row/column crossed as contributing double the distance (expansion factor 2). Sum all pairwise distances.

**Approach:** Parse galaxies, track empty rows/columns, calculate distances accounting for expansions without actually expanding the grid.

**Complexity:** O(n²) where n is the number of galaxies.

## Part 2
Same algorithm but with expansion factor 1,000,000 instead of 2.

**Approach:** Parameterize the expansion factor and reuse part 1 logic.

## Run
```bash
cargo run -q --bin day11 < input.txt
```

## Test
```bash
cargo test day11::
```
