# Day 4: Giant Squid

Mark each draw on every board. A full row or column wins. Keep the first score
for part 1 and the last score for part 2, skipping boards that already won.
Score = sum of unmarked numbers times the winning draw.

With fixed 5-by-5 boards, time is `O(draws * boards)` and space is `O(boards)`.

Boards are parsed as fixed 5-by-5 grids. Each draw marks matching cells on every
board; completed rows and columns identify winners. The solver records the first
winning score and continues until every board has won, returning the last score.
It rejects malformed boards, draws, games with no boards, and games where not
every board wins.

With `D` draws and `B` boards, runtime is `O(D * B)` because each board has a
fixed 25 cells. Space usage is `O(B)`.

From `2021/go`:

```sh
go run ./day04 < day04/testdata/example.txt
go test -race ./day04
go vet ./day04
```

The example prints exactly:

```text
Part 1: 4512
Part 2: 1924
```
