# Day 14: Docking Data

## Part A: Bitmask Application

Apply a bitmask to values before storing in memory. The mask contains:
- `0`: bit must be 0
- `1`: bit must be 1
- `X`: bit is left unchanged

Sum all memory values at the end.

**Complexity:** O(n) where n is the number of instructions.

## Part B: Floating Address Decoding

Apply a bitmask to memory addresses, where `X` creates floating bits:
- `0`: bit is unchanged from the address
- `1`: bit is set to 1
- `X`: bit "floats" (generates all combinations)

Each floating bit doubles the number of addresses written. The solution generates all 2^k addresses where k is the number of floating bits, capped at 20 to prevent enumerating more than ~1 million addresses.

Returns an error if a mask has more than 20 floating bits.

**Complexity:** O(n * 2^k) where n is instructions and k ≤ 20 is floating bits per instruction.

## Example

- Part A: 165
- Part B: 208

## Commands

```bash
cd 2020/go
go test -race ./challenge/day14
go run . 14 < ../rust/input/2020/day14-example.txt
go run . 14 < ../rust/input/2020/day14.txt
```
