# Day 18: Lavaduct Lagoon

Each instruction traces one edge of the trench polygon. The shoelace formula
gives the signed area enclosed by the trench's center line; Pick's theorem
(`area = interior + boundary/2 - 1`) then converts that into the count of
dug-out unit tiles, `interior + boundary = area + boundary/2 + 1`, so the
boundary length (the perimeter) never needs to be walked tile by tile.

Part 1 reads direction letters and distances directly. Part 2 decodes the
same geometry from the hex colour: the first five hex digits are the distance,
the last digit is the direction (`0`=R, `1`=D, `2`=L, `3`=U).

Both parts run in `O(n)` for `n` instructions and `O(1)` extra space.

From `2023/rust`:

```sh
cargo test day18::
cargo clippy --all-targets -- -D warnings
cargo run -q --bin day18 < examples/day18.txt
```

Expected output:

```text
Part 1: 62
Part 2: 952408144115
```
