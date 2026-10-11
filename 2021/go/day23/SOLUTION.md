# Day 23: Amphipod

Both parts run the same search: Dijkstra over burrow states, a state being the
hallway plus the occupied slots of the four rooms. Only two moves exist - the
amphipod nearest the hallway steps out of a room that still holds a stranger, or
an amphipod waiting in the hallway walks into its own room once that room holds
nothing but its own kind. A room-to-room walk is those two in sequence, stopping
on a free cell between the doors, so no cheaper route is missed. Walking an
amphipod home is always part of an optimal solution, so the generator takes that
move alone whenever it exists, which keeps the search small.

Part 2 unfolds the burrow: the rows `#D#C#B#A#` and `#D#B#A#C#` slide between the
two rows of the diagram, giving four slots per room and the same search.

For `S` reachable states and at most 28 moves each, time is `O(S log S)` and
space `O(S)`. The state is a comparable fixed-size struct, so it is its own key
in the distance map.

From `2021/go`:

```sh
go run ./day23 < day23/testdata/example.txt
go test -race ./day23
go vet ./day23
```

The example prints exactly:

```text
Part 1: 12521
Part 2: 44169
```
