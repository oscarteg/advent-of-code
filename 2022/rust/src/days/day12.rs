use anyhow::{bail, ensure, Context, Result};
use std::collections::VecDeque;

struct Grid {
    heights: Vec<u8>,
    width: usize,
    height: usize,
    start: usize,
    end: usize,
}

impl Grid {
    fn parse(input: &str) -> Result<Self> {
        let rows: Vec<&str> = input.lines().filter(|line| !line.is_empty()).collect();
        ensure!(!rows.is_empty(), "grid must have at least one row");
        let width = rows[0].len();
        ensure!(width > 0, "grid rows must not be empty");

        let mut heights = Vec::with_capacity(width * rows.len());
        let mut start = None;
        let mut end = None;
        for row in &rows {
            ensure!(
                row.len() == width,
                "all grid rows must have the same length"
            );
            for ch in row.bytes() {
                let index = heights.len();
                let elevation = match ch {
                    b'S' => {
                        ensure!(start.is_none(), "grid must have exactly one start");
                        start = Some(index);
                        b'a'
                    }
                    b'E' => {
                        ensure!(end.is_none(), "grid must have exactly one end");
                        end = Some(index);
                        b'z'
                    }
                    b'a'..=b'z' => ch,
                    other => bail!("invalid elevation character: {}", other as char),
                };
                heights.push(elevation - b'a');
            }
        }

        Ok(Self {
            heights,
            width,
            height: rows.len(),
            start: start.context("grid must have a start")?,
            end: end.context("grid must have an end")?,
        })
    }

    fn neighbors(&self, index: usize) -> impl Iterator<Item = usize> + '_ {
        let (row, col) = (index / self.width, index % self.width);
        [
            (row > 0).then(|| index - self.width),
            (row + 1 < self.height).then(|| index + self.width),
            (col > 0).then(|| index - 1),
            (col + 1 < self.width).then(|| index + 1),
        ]
        .into_iter()
        .flatten()
    }

    /// Shortest number of steps from any of `starts` to `self.end`, climbing
    /// at most one elevation per step, via breadth-first search from the end.
    fn shortest_path_from_end(&self, starts: &[usize]) -> Option<u32> {
        let mut distances = vec![u32::MAX; self.heights.len()];
        let mut queue = VecDeque::new();
        distances[self.end] = 0;
        queue.push_back(self.end);

        while let Some(current) = queue.pop_front() {
            let current_distance = distances[current];
            for neighbor in self.neighbors(current) {
                if distances[neighbor] != u32::MAX {
                    continue;
                }
                // Reversed edge: moving from `neighbor` to `current` in the
                // forward direction must drop by at most one elevation.
                if self.heights[current] as i16 - self.heights[neighbor] as i16 > 1 {
                    continue;
                }
                distances[neighbor] = current_distance + 1;
                queue.push_back(neighbor);
            }
        }

        starts
            .iter()
            .map(|&start| distances[start])
            .filter(|&distance| distance != u32::MAX)
            .min()
    }
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let grid = Grid::parse(input)?;
    let part1 = grid
        .shortest_path_from_end(&[grid.start])
        .context("no path from start to end")?;
    let lowest_points: Vec<usize> = grid
        .heights
        .iter()
        .enumerate()
        .filter(|(_, &elevation)| elevation == 0)
        .map(|(index, _)| index)
        .collect();
    let part2 = grid
        .shortest_path_from_end(&lowest_points)
        .context("no path from any lowest point to end")?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day12-example.txt");
        assert_eq!(solve(input).unwrap(), (31, 29));
    }

    #[test]
    fn single_row_staircase() {
        let input = "SbcdefghijklmnopqrstuvwxyE";
        assert_eq!(solve(input).unwrap(), (25, 25));
    }

    #[test]
    fn missing_start_is_rejected() {
        let input = "ab\nab\nabE";
        assert!(solve(input).is_err());
    }

    #[test]
    fn ragged_rows_are_rejected() {
        let input = "Sab\nabE";
        assert!(solve(input).is_err());
    }

    #[test]
    fn invalid_character_is_rejected() {
        let input = "S1\naE";
        assert!(solve(input).is_err());
    }
}
