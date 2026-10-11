# Day 19: Aplenty

- **Part 1:** Parse each workflow into rules of `category`, `<`/`>`,
  threshold and target. Walk every part from `in`, always following the first
  rule whose condition holds, until it reaches `A` or `R`, then sum the four
  ratings of the accepted parts.
- **Part 2:** Push the whole search space (`1..=4000` per category) through the
  same rules instead of a single part. Each condition splits the incoming space
  into the ratings it routes away and the ratings the next rule sees, so the
  accepted fragments are disjoint: the answer is the sum of the products of
  their four span widths.
- Part 1 is `O(parts * rules)`. Part 2 visits every rule once per reachable
  space fragment and never enumerates ratings. Parsing rejects malformed input
  (unknown category, missing fallback rule, unknown target, cyclic workflows).

From `2023/rust`:

```sh
cargo run -q --bin day19 < examples/day19.txt # 19114
cargo run -q --bin day19 < examples/day19.txt # 167409079868000
cargo test day19::
```

Use the same commands with your puzzle input redirected into stdin.
