# Day 1: Trebuchet?!

## Part 1

Extract the first and last ASCII digit from each line, concatenate them, and sum.

- Iterate through each character in the line
- Collect all ASCII digits
- Pair the first and last digits (use same digit twice if only one found)
- Parse as two-digit number and add to sum

**Complexity:** O(n) where n is total input length.

## Part 2

Same process, but also recognize spelled-out digits ("one", "two", ..., "nine").

- Iterate through each position in the line
- Check if a digit word starts at that position, or if the character is an ASCII digit
- Collect all found digits and pair first + last

Overlapping digits (e.g., "eightwo") are handled by checking every position.

**Complexity:** O(n × 9) = O(n).

## Run and Test

```bash
cd 2023/go

# Test
go test -race ./day01

# Run
go run ./day01 < day01/testdata/example.txt    # part 1 example: 142
go run ./day01 < day01/testdata/example2.txt   # part 2 example: 281
```
