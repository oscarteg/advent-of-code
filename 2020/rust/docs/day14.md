# Day 14: Docking Data

## Part 1: Bitmask Application

Apply a bitmask to values before storing in memory. The mask contains:
- `0`: bit must be 0
- `1`: bit must be 1  
- `X`: bit is left unchanged

Sum all memory values at the end.

**Complexity:** O(n) where n is the number of instructions.

## Part 2: Floating Address Decoding

Apply a bitmask to memory addresses, where `X` creates floating bits:
- `0`: bit is unchanged from the address
- `1`: bit is set to 1
- `X`: bit "floats" (generates all combinations)

Each floating bit doubles the number of addresses written. The solution generates all 2^k addresses where k is the number of floating bits, capped at 20 to avoid enumerating too many addresses (>1 million combinations).

Returns an error if a mask has more than 20 floating bits.

**Complexity:** O(n * 2^k) where n is instructions and k ≤ 20 is floating bits per instruction.

## Commands

```bash
cargo test --bin day14
cargo run -q --bin day14 < input/2020/day14-example.txt
```

## Example

**Part 1:** 165 (from inline test)  
**Part 2:** 208 (from day14-example.txt)
