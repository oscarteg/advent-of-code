# Day 16: The Floor Will Be Lava

## Approach

**Part 1:** Simulate a laser beam entering the grid at the top-left corner moving right. Track which tiles the beam energizes using BFS. Use `(position, direction)` pairs in a `HashSet` to detect cycles and avoid infinite loops.

**Part 2:** Try all possible starting positions (all edges with appropriate entry directions) and find which configuration energizes the most tiles.

## Algorithm

- Represent the grid as a 2D vector of characters
- Use isize coordinates to handle boundary checks cleanly (out-of-bounds values fail the bounds check naturally)
- For each cell type, apply the appropriate transformation:
  - `.`: continue in the same direction
  - `/`: reflect using direction swaps and negations
  - `\`: reflect using different swaps
  - `-`: if horizontal, continue; if vertical, split into left and right
  - `|`: if vertical, continue; if horizontal, split into up and down
- Track visited (position, direction) pairs to detect cycles
- Count unique positions in the visited set

## Complexity

- Time: O(rows × cols × 4) for part 1; O(edges × rows × cols × 4) for part 2
- Space: O(rows × cols) for visited positions

## Commands

```bash
cargo test day16::
cargo clippy --all-targets -- -D warnings
cargo run -q --bin day16 < examples/day16.txt
```
