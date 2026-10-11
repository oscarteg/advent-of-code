use std::cmp::Ordering;
use std::collections::BinaryHeap;

use anyhow::{bail, ensure, Context, Result};

struct Grid {
    risks: Vec<u8>,
    width: usize,
    height: usize,
}

impl Grid {
    fn parse(input: &str) -> Result<Self> {
        let rows: Vec<&str> = input.lines().collect();
        ensure!(!rows.is_empty(), "missing grid rows");
        let width = rows[0].len();
        ensure!(width > 0, "grid rows must not be empty");
        let mut risks = Vec::with_capacity(width * rows.len());
        for row in &rows {
            ensure!(row.len() == width, "grid rows must have equal length");
            for ch in row.chars() {
                let digit = ch
                    .to_digit(10)
                    .with_context(|| format!("invalid risk digit '{ch}'"))?;
                risks.push(u8::try_from(digit).context("risk digit out of range")?);
            }
        }
        Ok(Self {
            risks,
            width,
            height: rows.len(),
        })
    }

    fn risk_at(&self, x: usize, y: usize) -> u32 {
        u32::from(self.risks[y * self.width + x])
    }

    /// Builds the tiled grid: `tiles`x`tiles` copies, each increasing risk by
    /// the Manhattan tile distance, wrapping back to 1 after 9.
    fn tile(&self, tiles: usize) -> TiledGrid<'_> {
        TiledGrid {
            base: self,
            tiles,
            width: self.width * tiles,
            height: self.height * tiles,
        }
    }
}

struct TiledGrid<'a> {
    base: &'a Grid,
    tiles: usize,
    width: usize,
    height: usize,
}

impl TiledGrid<'_> {
    fn risk_at(&self, x: usize, y: usize) -> u32 {
        let base_risk = self.base.risk_at(x % self.base.width, y % self.base.height);
        let tile_distance = (x / self.base.width + y / self.base.height) as u32;
        (base_risk - 1 + tile_distance) % 9 + 1
    }

    fn lowest_total_risk(&self) -> Result<u32> {
        ensure!(self.tiles > 0, "tile count must be positive");
        dijkstra(self.width, self.height, |x, y| self.risk_at(x, y))
    }
}

#[derive(Eq, PartialEq)]
struct State {
    risk: u32,
    x: usize,
    y: usize,
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> Ordering {
        other.risk.cmp(&self.risk)
    }
}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

fn dijkstra(width: usize, height: usize, risk_at: impl Fn(usize, usize) -> u32) -> Result<u32> {
    let mut best = vec![u32::MAX; width * height];
    let mut heap = BinaryHeap::new();
    best[0] = 0;
    heap.push(State {
        risk: 0,
        x: 0,
        y: 0,
    });

    let target = (width - 1, height - 1);
    while let Some(State { risk, x, y }) = heap.pop() {
        if (x, y) == target {
            return Ok(risk);
        }
        if risk > best[y * width + x] {
            continue;
        }
        let mut neighbors = Vec::with_capacity(4);
        if x > 0 {
            neighbors.push((x - 1, y));
        }
        if x + 1 < width {
            neighbors.push((x + 1, y));
        }
        if y > 0 {
            neighbors.push((x, y - 1));
        }
        if y + 1 < height {
            neighbors.push((x, y + 1));
        }
        for (nx, ny) in neighbors {
            let next_risk = risk + risk_at(nx, ny);
            let index = ny * width + nx;
            if next_risk < best[index] {
                best[index] = next_risk;
                heap.push(State {
                    risk: next_risk,
                    x: nx,
                    y: ny,
                });
            }
        }
    }
    bail!("no path to the bottom-right corner")
}

fn solve_with(input: &str, tiles: usize) -> Result<(u32, u32)> {
    let grid = Grid::parse(input)?;
    let part1 = grid.tile(1).lowest_total_risk()?;
    let part2 = grid.tile(tiles).lowest_total_risk()?;
    Ok((part1, part2))
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    solve_with(input, 5)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day15-example.txt");
        assert_eq!(solve(input).unwrap(), (40, 315));
    }

    #[test]
    fn single_cell_grid_has_zero_risk() {
        assert_eq!(solve_with("5", 1).unwrap(), (0, 0));
    }

    #[test]
    fn risk_wraps_after_nine_when_tiled() {
        let input = include_str!("../../input/day15-example.txt");
        // A tile count of 1 must match the untiled grid's risk.
        assert_eq!(solve_with(input, 1).unwrap(), (40, 40));
    }

    #[test]
    fn uneven_rows_are_rejected() {
        assert!(solve("123\n12").is_err());
    }

    #[test]
    fn non_digit_cells_are_rejected() {
        assert!(solve("12a\n345").is_err());
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("").is_err());
    }
}
