# Day 21: Keypad Conundrum (Zig)

Every code is typed on the door's numeric keypad by a robot that is itself commanded from a
directional keypad by another robot, and so on; the human types the last keypad. Part 1 chains
two directional-keypad robots and part 2 chains 25, and both sum each code's complexity: the
presses the human needs times the code's numeric part.

## Approach

A robot arm only walks a shortest path between two keys, and just two orders can be shortest:
every horizontal step first, or every vertical step first - one of them possibly blocked
because its turning point is the keypad's gap. Pressing a key therefore costs the cheaper of
typing `route + "A"` one keypad up, and with no robot left the human simply presses the route's
buttons. A move's cost depends on nothing but the keypad, the key pair and the robots above it,
so memoising on those three collapses the exponential expansion: every sequence length is a sum
of cached move costs. That memo holds `robots * keys^2` entries - 25 key pairs per directional
level - which is also the time bound per code.

## Run / test

```bash
cd 2024/zig
zig test day21.zig
zig run day21.zig < examples/day21.txt
```
