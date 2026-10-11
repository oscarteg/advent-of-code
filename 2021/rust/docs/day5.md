# Day 5: Hydrothermal Venture

## Problem

Count overlapping hydrothermal vent lines on a 2D grid. Lines are represented as segments from `(x1,y1)` to `(x2,y2)`.

- **Part 1**: Count points where 2+ horizontal or vertical lines overlap
- **Part 2**: Count points where 2+ lines (including diagonals) overlap

## Algorithm

1. **Parse** each line segment into (x1, y1, x2, y2) coordinates
2. **Categorize** lines:
   - Horizontal: y1 == y2
   - Vertical: x1 == x2
   - Diagonal: |x2 - x1| == |y2 - y1| (45° only)
3. **Mark grid points** using direction vectors (signum) to step through each line
4. **Count** points with 2+ line passes for each part

## Complexity

- **Time**: O(n × m) where n = number of lines, m = max line length (up to 1000)
- **Space**: O(p) where p = unique points covered (max 1M for 1000×1000 grid)

## Implementation Details

- `Line::points_on_line()` uses direction vectors to walk both horizontal and diagonal lines uniformly
- Two separate HashMaps track overlaps: one for H/V only (part 1), one for all lines (part 2)
- Early filtering of line types avoids processing irrelevant lines

## Running

```bash
cd 2021/rust

cargo test --bin day5              # Run tests
cargo clippy --bin day5 -- -D warnings  # Lint
cargo run -q --bin day5 < input/day5-example.txt  # CLI with example
```

Expected output on example:
```
Part 1: 5
Part 2: 12
```

## Tests

- **official_example**: Validates both parts on provided puzzle input
- **horizontal_line_overlap**: Overlapping horizontal segments
- **vertical_line_overlap**: Overlapping vertical segments
- **diagonal_lines_part2_only**: Diagonal lines ignored in part 1
- **single_point_overlap**: Single overlapping point counts
- **non_axis_aligned_non_diagonal_ignored**: Arbitrary angle lines filtered
- **malformed_input**: Missing arrow and invalid coordinates rejected
- **empty_input**: Empty input produces 0 overlaps
- **empty_lines_ignored**: Whitespace-only lines skip cleanly
