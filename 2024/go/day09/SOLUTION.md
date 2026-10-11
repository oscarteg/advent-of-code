# Day 9: Disk Fragmenter

Expand the dense disk map into segments `{id, start, length}` (`id = -1` for
free space). Part 1 fragments files block by block with a two-pointer
compaction over an explicit block array. Part 2 moves whole files, highest
id first, into the leftmost free span that fits, without fragmenting them.

## Complexity

Part 1: O(n) in the number of blocks. Part 2: O(files x free spans).

## Run / test

```bash
cd 2024/go && go test -race -count=1 ./day09
go run ./day09 < day09/testdata/example.txt
```
