# Day 11: Monkey in the Middle

## Part 1
Parse each monkey's starting items, worry operation, divisor and throw
targets. Simulate 20 rounds: every monkey inspects each item (applying its
operation), divides worry by 3, then throws it based on divisibility.
Answer is the product of the two highest inspection counts.
Complexity: O(rounds * items).

## Part 2
Dropping the divide-by-3 relief makes worry grow unbounded. Since only
divisibility by each monkey's divisor matters, worry is instead reduced
modulo the product of all divisors after every inspection, keeping values
small while preserving every divisibility test. Run for 10,000 rounds.
Complexity: O(rounds * items).

## Commands
```bash
go test -race -count=1 ./day11
go run ./day11 < day11/testdata/example.txt
```
