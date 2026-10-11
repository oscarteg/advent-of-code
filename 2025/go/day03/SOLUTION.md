# Day 3: Lobby

## Approach

Each line represents a bank of digit batteries. We pick K batteries (digits) from each bank while preserving their order in the string, forming the largest possible K-digit number.

**Algorithm:** Greedy per position. For each digit position in the result:
- Search remaining input for the largest digit
- Ensure enough digits remain after it for remaining positions
- Use that digit and continue from the next position

Part 1: pick 2 digits per bank, sum all. Part 2: pick 12 digits per bank, sum all.

## Complexity

- Time: O(n * k²) where n is total digits and k is target length (k=2 or k=12)
- Space: O(k) for result

## Run

```bash
cd 2025/go
go run ./day03 < day03/testdata/example.txt
go test -race -count=1 ./day03
```
