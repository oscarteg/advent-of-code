# Day 17: Trick Shot

## Approach

**Part 1**: Find the initial velocity that reaches the highest y-coordinate while still landing in the target area. For a given vx and vy, simulate the probe motion step by step until it either hits the target or overshoots. Track the maximum y reached for all valid velocities.

**Part 2**: Count the number of distinct initial velocities (vx, vy) that result in the probe landing in the target area.

## Physics

The probe follows simple physics:
- Each step: position += velocity, then velocity -= (1, 1) for gravity
- Horizontal velocity stops decreasing once it reaches 0
- We need to find all velocities where the probe intersects the target rectangle

## Bounds

- vx ranges from 1 to target.x_max (any higher and we overshoot horizontally)
- vy ranges from target.y_min - 1 to roughly 2*|target.y_min| (symmetric for upward and downward trajectories)

## Complexity

O(vx_range × vy_range × simulation_steps) per solution. With reasonable bounds, this is fast enough.

## Commands

```bash
cd 2021/rust
cargo test --bin day17
cargo run -q --bin day17 < input/day17-example.txt
```
