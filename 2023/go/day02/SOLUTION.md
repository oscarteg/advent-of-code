# Day 2: Cube Conundrum

Each game is a semicolon-separated list of reveals (red/green/blue cube
counts). `parseGames` turns every line into a `game` with its id and reveals.

**Part 1**: sum the ids of games where every reveal fits within 12 red, 13
green, 14 blue cubes.

**Part 2**: for each game, take the maximum of each color across all reveals
(the fewest cubes that could have produced it), multiply the three maxima
(the "power"), and sum the powers.

Both parts share one pass over the parsed games: `O(games * reveals)` time
and `O(games)` space.

From `2023/go`:

```sh
go test -race ./day02
go vet ./day02
go run ./day02 < day02/testdata/example.txt
```

The example prints:

```text
Part 1: 8
Part 2: 2286
```
