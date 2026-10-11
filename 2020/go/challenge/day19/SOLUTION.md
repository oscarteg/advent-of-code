# Day 19: Monster Messages

## Approach

**Part A:** Parse the rule block into a `map[int]rule`, where a rule is either a literal byte or
alternative sequences of rule IDs. Match rule 0 against a message by tracking *every* position a
rule can end at: a literal advances one byte, a sequence threads the position set through its
sub-rules, alternatives union their results. The message matches when one of rule 0's end
positions equals its length.

**Part B:** Replace rule 8 with `42 | 42 8` and rule 11 with `42 31 | 42 11 31`, then count again.
Position sets already handle the recursion, because both loops match rule 42 before recursing and
therefore always consume input. Rule sets without rules 8/11 are left alone, so the part A example
keeps its answer. End positions are memoised per (rule, start) while matching one message.

## Complexity

With `m` messages of length `n` and `r` rules: O(m · r · n²), since the cache holds one set of at
most `n` positions per (rule, start) pair.

## Commands

```bash
cd 2020/go
go test -race ./challenge/day19
go run . 19 < ../rust/input/2020/day19.txt   # once main.go registers day 19
```
