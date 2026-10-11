# Day 16: The Floor Will Be Lava

## Approach

**Part 1:** Simulate a laser beam entering the grid at the top-left corner moving right. Use BFS to track the beam's path, with visited `(row, col, direction)` triples to detect cycles.

**Part 2:** Try all possible starting positions from the four edges and find which one energizes the most tiles.

## Algorithm

- Parse the grid into a 2D slice of runes
- Use a queue-based BFS for simulation with state `(row, col, direction)`
- Track visited `(row, col, direction)` tuples to avoid infinite loops
- For each cell type, apply appropriate transformations:
  - `.`: continue in the same direction
  - `/`: reflect using directional mappings
  - `\`: reflect using different directional mappings
  - `-`: pass through if horizontal, split if vertical
  - `|`: pass through if vertical, split if horizontal
- Count unique `(row, col)` positions in the visited set

## Complexity

- Time: O(rows × cols × 4) for part 1; O(edges × rows × cols × 4) for part 2
- Space: O(rows × cols) for visited positions

## Commands

```bash
go test ./day16
go vet ./day16
go run ./day16 < day16/testdata/example.txt
```
