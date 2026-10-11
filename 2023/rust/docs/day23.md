# Day 23: A Long Walk

Parsing turns the map into tiles, the entrance (the only path tile of the top
row) and the exit (the only one of the bottom row). The trails are corridors
between forks, so both parts collapse the grid: nodes are the entrance, the exit
and every tile with more than two walkable neighbours, and following each
corridor from a node to the next node yields a weighted edge. Finding the
longest path is NP-hard, so the search enumerates every hike that visits each
node at most once and keeps the longest.

- **Part 1:** slopes are icy. A step onto a slope forces the next step downhill,
  which makes the edges directed and prunes corridors that end in a slope
  pointing back.
- **Part 2:** slopes count as paths, so the edges are symmetric and far more
  hikes exist.

Compressing costs one pass over the grid; the search is exponential in the
number of nodes (about 36 on a real input) but linear per enumerated hike.

From `2023/rust`:

```sh
cargo run -q --bin day23 < examples/day23.txt # 94
cargo run -q --bin day23 < examples/day23.txt # 154
cargo test day23::
```
