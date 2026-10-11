# Day 13: Distress Signal

## Approach

**Part 1:** Parse packet pairs from input, compare each pair using custom ordering rules:
- Integer vs integer: compare numerically
- List vs list: compare element-by-element; shorter list is "less"
- Integer vs list: wrap integer in list and compare

Sum the 1-indexed pair indices where left < right (correct order).

**Part 2:** Parse all packets (skip blank lines), add divider packets `[[2]]` and `[[6]]`, sort using the same comparison, find the 1-indexed positions of both dividers, return their product.

## Complexity

- **Time:** O(n log n) for sorting in part 2; O(n * d) for comparisons where d is max packet depth
- **Space:** O(n) for packet storage

## Run & Test

```bash
cd 2022/go
test -z "$(gofmt -l day13)" && go vet ./day13 && go test -race -count=1 ./day13
go run ./day13 < day13/testdata/example.txt
```

Expected output: `Part 1: 13` and `Part 2: 140`
