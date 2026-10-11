use crate::day21::custom_error::AocError;
use std::collections::HashSet;

type Plot = (i64, i64);

#[derive(Clone, Copy)]
enum Tiling {
    Single,
    Infinite,
}

#[derive(Debug)]
pub(crate) struct Garden {
    /// Row-major, `width * height` long; `true` is a rock.
    rocks: Vec<bool>,
    width: i64,
    height: i64,
    start: Plot,
}

impl Garden {
    pub fn parse(input: &str) -> Result<Self, AocError> {
        let lines: Vec<&str> = input
            .lines()
            .map(str::trim)
            .filter(|line| !line.is_empty())
            .collect();
        let width = lines
            .first()
            .ok_or(AocError::InvalidInput("map is empty"))?
            .len();
        let mut rocks = Vec::with_capacity(width * lines.len());
        let mut start = None;
        for (row, line) in lines.iter().enumerate() {
            if line.len() != width {
                return Err(AocError::InvalidInput("rows differ in length"));
            }
            for (column, tile) in line.bytes().enumerate() {
                match tile {
                    b'.' => rocks.push(false),
                    b'#' => rocks.push(true),
                    b'S' if start.is_some() => {
                        return Err(AocError::InvalidInput("more than one start"))
                    }
                    b'S' => {
                        start = Some((column as i64, row as i64));
                        rocks.push(false);
                    }
                    _ => return Err(AocError::InvalidInput("unknown tile")),
                }
            }
        }
        Ok(Self {
            rocks,
            width: width as i64,
            height: lines.len() as i64,
            start: start.ok_or(AocError::InvalidInput("no start position"))?,
        })
    }

    /// Plots reachable in exactly `steps` steps, the map edges acting as rocks.
    pub fn reachable(&self, steps: u64) -> u64 {
        self.flood(steps, Tiling::Single)
    }

    /// Plots reachable in exactly `steps` steps on the infinitely tiled map.
    pub fn reachable_tiled(&self, steps: u64) -> u64 {
        self.flood(steps, Tiling::Infinite)
    }

    /// Plots reachable on the tiled map for step counts too large to flood.
    ///
    /// Every `width` steps the flood swallows one more copy of the map in each
    /// direction, so on a square map with an open border and an open row and
    /// column through the start the counts after `offset + tiles * width` steps
    /// are a quadratic in `tiles`. Fit that quadratic through the first three
    /// tile-aligned floods and evaluate it at the requested step count.
    pub fn reachable_extrapolated(&self, steps: u64) -> Result<u64, AocError> {
        let width = self.width as u64;
        if steps < 2 * width {
            return Ok(self.reachable_tiled(steps));
        }
        if self.width != self.height {
            return Err(AocError::InvalidInput("extrapolation needs a square map"));
        }
        let offset = steps % width;
        let [first, second, third] =
            [0, 1, 2].map(|tiles| i128::from(self.reachable_tiled(offset + tiles * width)));
        let tiles = i128::from(steps / width);
        let plots = first
            + tiles * (second - first)
            + tiles * (tiles - 1) / 2 * (third - 2 * second + first);
        u64::try_from(plots).map_err(|_| AocError::InvalidInput("counts are not quadratic"))
    }

    fn flood(&self, steps: u64, tiling: Tiling) -> u64 {
        // The elf can step off a plot and back onto it, so every plot whose
        // distance has the same parity as the step count is reachable - unless
        // he is walled in and cannot take a single step.
        let can_step = neighbours(self.start)
            .iter()
            .any(|&plot| self.is_plot(plot, tiling));
        if steps > 0 && !can_step {
            return 0;
        }
        let mut seen = HashSet::from([self.start]);
        let mut frontier = vec![self.start];
        let mut reachable = 0;
        for remaining in (0..=steps).rev() {
            if remaining % 2 == 0 {
                reachable += frontier.len() as u64;
            }
            let next: Vec<Plot> = frontier
                .iter()
                .flat_map(|&plot| neighbours(plot))
                .filter(|&plot| self.is_plot(plot, tiling) && seen.insert(plot))
                .collect();
            frontier = next;
        }
        reachable
    }

    fn is_plot(&self, (x, y): Plot, tiling: Tiling) -> bool {
        let (column, row) = match tiling {
            Tiling::Infinite => (x.rem_euclid(self.width), y.rem_euclid(self.height)),
            Tiling::Single if (0..self.width).contains(&x) && (0..self.height).contains(&y) => {
                (x, y)
            }
            Tiling::Single => return false,
        };
        !self.rocks[(row * self.width + column) as usize]
    }
}

fn neighbours((x, y): Plot) -> [Plot; 4] {
    [(x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)]
}
