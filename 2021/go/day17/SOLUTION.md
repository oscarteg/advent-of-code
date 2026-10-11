# Day 17: Trick Shot

## Approach

**Part 1**: Simulate all reasonable initial velocities (vx, vy) and find the one that reaches the highest y-coordinate while landing in the target area. Track max y for all valid trajectories.

**Part 2**: Count the total number of distinct velocities that result in hitting the target.

## Physics

Each step: position advances by velocity, then velocity decreases by 1 in both dimensions (gravity).
- Horizontal velocity clamps at 0 (stops decelerating)
- Continue simulation until probe either hits target or overshoots

## Bounds

- vx: 1 to target.xMax (going higher overshoots immediately)
- vy: target.yMin - 1 to 2×|target.yMin| (symmetric for up/down trajectories)

## Complexity

O(vx_range × vy_range × simulation_steps). With typical puzzle sizes, negligible.

## Commands

```bash
cd 2021/go
go test ./day17
go run ./day17 < day17/testdata/example.txt
```
