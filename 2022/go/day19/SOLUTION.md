# Day 19: Not Enough Minerals

## Part 1
Parse every blueprint into robot costs per resource, then search over *which
robot to build next*: wait the minutes a robot's missing resources need at the
current collection rate, spend one minute building it, repeat. A geode robot
finished with `m` minutes left cracks exactly `m` geodes, so its yield is
banked at once and geode robots need no further tracking. Part 1 sums
`id * max geodes in 24 minutes` over all blueprints.

## Part 2
The same search over 32 minutes for the first three blueprints, multiplied.

## Pruning
- Never build more ore, clay or obsidian robots than the largest cost of that
  resource in any recipe; the surplus can never be spent.
- Skip a robot that cannot be finished with a minute left to collect.
- Bound a branch with a relaxed game that has free ore and allows one robot of
  every type per minute; drop the branch when that bound cannot beat the best
  schedule found so far. The worst case stays exponential, the example is
  instant.

## Run / test
```sh
go run ./day19 < day19/testdata/example.txt
go test -race ./day19
```

The example prints `Part 1: 33` and `Part 2: 3472`.
