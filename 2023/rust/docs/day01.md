# Day 1: Trebuchet?!

## Part 1

Extract the first and last digit from each line, concatenate them, and sum.

- Find all ASCII digits in each line
- Pair the first and last; if only one digit, use it twice
- Parse the pair as a two-digit number
- Sum across all lines

**Complexity:** O(n) where n is the total input length.

## Part 2

Same process, but also recognize spelled-out digits ("one", "two", ..., "nine").

- Iterate through each position in the line
- Check if a digit word starts at that position, or if the character is an ASCII digit
- Collect all found digits (both forms) and pair first + last

Overlapping digits (e.g., "eightwo") are handled by checking every position.

**Complexity:** O(n × 9) = O(n) for checking each of 9 digit words at each position.

## Run and Test

```bash
cd 2023/rust

# Test
cargo test day01::

# Run
cargo run -q --bin day01 < examples/day01.txt
cargo run -q --bin day01 < examples/day01-2.txt
```
