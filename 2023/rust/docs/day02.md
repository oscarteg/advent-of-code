# Day 2: Cube Conundrum

Each game is a list of reveals (red/green/blue cube counts). Shared parsing in
`cube_game.rs` turns a line like `Game 1: 3 blue, 4 red; ...` into a `Game`
with its id and reveals.

**Part 1**: sum the ids of games where every reveal fits within 12 red, 13
green, 14 blue cubes.

**Part 2**: for each game, take the maximum of each color across all reveals
(the fewest cubes that could have produced it), multiply the three maxima
(the "power"), and sum the powers.

Both parts are `O(games * reveals)` time and `O(games)` space.

From `2023/rust`:

```sh
cargo test day02::
cargo clippy --all-targets -- -D warnings
cargo run -q --bin day02 < examples/day02.txt
```

The example prints `8` for part 1 and `2286` for part 2.
