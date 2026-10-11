# Day 15: Lens Library

Part 1: run the HASH algorithm (multiply by 17, mod 256 per byte) over each
comma-separated step and sum the results. `O(n)` in input length.

Part 2: parse each step as `label-` (remove) or `label=focal` (insert), hash
the label to pick one of 256 boxes, then insert/update/remove the lens in
that box's ordered list. Sum `(box + 1) * (slot + 1) * focal_length` over all
boxes. `O(n * k)` where `k` is the average lenses per box (small in
practice), `O(n)` space.

From `2023/rust`:

```sh
cargo test day15::
cargo clippy --all-targets -- -D warnings
cargo run -q --bin day15 < examples/day15.txt
```

Example answers: Part 1 is `1320`, Part 2 is `145`.
