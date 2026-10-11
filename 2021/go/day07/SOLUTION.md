# Day 7: The Treachery of Whales

Part 1 minimizes total absolute distance, so an input median is an optimal
alignment point. Part 2 uses triangular fuel cost `d * (d + 1) / 2`; its
convex total is minimized at an integer adjacent to the mean, so the solver
checks the floor and ceiling.

Sorting costs O(n log n); each fuel calculation scans the positions in O(n).
Run from `2021/go`:

```sh
go test -race ./day07
go vet ./day07
go run ./day07 < day07/testdata/example.txt
```

The example prints `Part 1: 37` and `Part 2: 168`.
