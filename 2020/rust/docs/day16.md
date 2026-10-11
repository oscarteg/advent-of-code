# Day 16: Ticket Translation

## Part 1

Parse the notes into field rules (a name plus inclusive value ranges), your own
ticket and the nearby tickets. The scanning error rate is the sum of every
nearby value that no rule accepts: O(tickets * fields * ranges).

## Part 2

Drop the nearby tickets that contain such a value. Every remaining ticket
constrains each column to the fields that accept its value there, so start each
column with all matching fields, then repeatedly take the column that has a
single candidate left and strike that field from the other columns. Resolving
the order costs O(fields^2) eliminations on top of the O(tickets * fields^2)
candidate scan; the answer multiplies the values on your own ticket whose field
name starts with `departure`. Notes without such a field (both examples) leave
an empty product of 1.

Ambiguous notes, malformed ranges, missing sections and tickets with the wrong
number of values are reported as errors.

## Run & test

```bash
cargo test --bin day16
cargo run -q --bin day16 < input/2020/day16-example.txt
```
