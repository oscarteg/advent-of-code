# Day 16: Packet Decoder

The transmission is one hexadecimal string. A bit cursor over its nibbles reads
every packet header: a 3-bit version and a 3-bit type ID. Type ID 4 is a literal
value stored as five-bit groups; every other type ID is an operator whose
sub-packets are delimited either by a 15-bit length in bits or an 11-bit packet
count. Parsing builds a `Packet` tree once, and only the zero bits that pad the
final nibble may follow the outermost packet.

Part 1 adds the versions of every node of that tree.

Part 2 evaluates the tree: sum, product, minimum, maximum, greater than, less
than and equal to. Arithmetic is checked, so an overflow is an error.

For a transmission of `N` bits, both parts take `O(N)` time and `O(N)` space;
the tree holds at most one packet per 11 bits.

From `2021/rust`:

```sh
cargo run -q --bin day16 < input/day16-example.txt
# Part 1: 6
# Part 2: 2021
cargo test --bin day16
```

The fixture is the puzzle's worked literal example `D2FE28` (version 6, value
2021); the version-sum and operator examples live in the tests.
