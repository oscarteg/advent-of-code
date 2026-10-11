use crate::day17::custom_error::AocError;
use crate::day17::grid::{Direction, Grid};
use std::cmp::Ordering;
use std::collections::BinaryHeap;

#[derive(Clone, Copy, PartialEq, Eq)]
struct State {
    heat_loss: u32,
    x: usize,
    y: usize,
    direction: Option<Direction>,
    run: u32,
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> Ordering {
        other.heat_loss.cmp(&self.heat_loss)
    }
}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

/// Dijkstra over (position, last direction, run length), bounded by `min_run` and `max_run`
/// consecutive blocks in the same direction before the crucible may stop or turn.
pub(crate) fn least_heat_loss(grid: &Grid, min_run: u32, max_run: u32) -> Option<u32> {
    let goal = (grid.width() - 1, grid.height() - 1);
    let mut best: Vec<u32> =
        vec![u32::MAX; grid.width() * grid.height() * 4 * (max_run as usize + 1)];
    let index = |x: usize, y: usize, direction: usize, run: u32| {
        ((y * grid.width() + x) * 4 + direction) * (max_run as usize + 1) + run as usize
    };

    let mut heap = BinaryHeap::new();
    heap.push(State {
        heat_loss: 0,
        x: 0,
        y: 0,
        direction: None,
        run: 0,
    });

    while let Some(state) = heap.pop() {
        if (state.x, state.y) == goal && state.run >= min_run {
            return Some(state.heat_loss);
        }

        for direction in Direction::ALL {
            if state.direction == Some(direction.opposite()) {
                continue;
            }
            let continuing = state.direction == Some(direction);
            if continuing && state.run >= max_run {
                continue;
            }
            if !continuing && state.direction.is_some() && state.run < min_run {
                continue;
            }
            let Some((nx, ny)) = grid.step(state.x, state.y, direction) else {
                continue;
            };

            let run = if continuing { state.run + 1 } else { 1 };
            let heat_loss = state.heat_loss + grid.cost(nx, ny);
            let slot = index(nx, ny, direction as usize, run);
            if heat_loss < best[slot] {
                best[slot] = heat_loss;
                heap.push(State {
                    heat_loss,
                    x: nx,
                    y: ny,
                    direction: Some(direction),
                    run,
                });
            }
        }
    }

    None
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let grid = Grid::parse(input)?;
    let heat_loss =
        least_heat_loss(&grid, 0, 3).ok_or(AocError::InvalidInput("no path to the factory"))?;
    Ok(heat_loss.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day17.txt");
        assert_eq!("102", process(input)?);
        Ok(())
    }

    #[test]
    fn single_block_grid_needs_no_moves() {
        let input = "5";
        assert_eq!(process(input).unwrap(), "0");
    }

    #[test]
    fn straight_line_longer_than_three_must_detour() {
        let input = "1111\n9991\n1991\n1119";
        // The top row alone would need four straight moves, which is forbidden, so the
        // crucible must turn into the costlier middle rows before reaching the corner.
        assert_eq!(process(input).unwrap(), "14");
    }

    #[test]
    fn malformed_input_is_rejected() {
        for input in ["", "12\n1", "1a2\n345"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
