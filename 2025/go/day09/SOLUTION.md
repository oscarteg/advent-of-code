# Day 9: Movie Theater

Red tiles are listed as `x,y` in the order they are walked around the room.

**Part 1** takes any two red tiles as opposite corners of a rectangle whose
tiles count inclusively (`(|dx| + 1) * (|dy| + 1)`): scan all pairs, keep the largest.

**Part 2** joins consecutive red tiles with straight green lines and fills
the enclosed area green, so a rectangle counts only when every tile it covers
is red or green. All those lines sit on coordinates of the loop, so x and y
compress to the distinct coordinates plus one representative per gap, which
makes every strip of that grid a single colour. A sweep per row marks the
strips a line runs through and counts the vertical lines right of each strip
(an odd count means enclosed); a prefix sum over the uncovered strips then
checks a rectangle in four lookups. Both parts run in `O(n^2)` time and memory.

```sh
cd 2025/go
go test -race ./day09
go run ./day09 < day09/testdata/example.txt   # Part 1: 50, Part 2: 24
```

Part 2's statement needs a login, so its rule (same two red corners, only red
and green tiles) comes from the puzzle notes: it is the only reading that
produces the documented example answer of 24.
