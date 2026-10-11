# Day 16: Ticket Translation

## Part A
Parse the notes into field rules (a name plus inclusive ranges), your own ticket
and the nearby tickets. The scanning error rate sums every nearby value that no
rule accepts: O(tickets * fields * ranges).

## Part B
Drop the nearby tickets holding such a value. Each remaining ticket constrains a
column to the fields that accept its value there, so every column starts with all
matching fields; then repeatedly take the column with a single candidate left and
strike that field from the other columns. That is O(fields^2) eliminations on top
of the O(tickets * fields^2) candidate scan. The answer multiplies the values on
your own ticket whose field name starts with `departure`; notes without such a
field (both examples) leave an empty product of 1.

Ambiguous notes, malformed ranges, missing sections and tickets with the wrong
number of values are returned as errors.

## Run and test
```bash
go run . 16 < ../rust/input/2020/day16.txt
go test -race ./challenge/day16
```
