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
zig test src/day_11.zig
zig run src/day_11.zig < src/examples/day_11.txt
```
