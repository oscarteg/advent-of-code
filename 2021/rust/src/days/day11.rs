use anyhow::{ensure, Context, Result};

const SIZE: usize = 10;

struct Grid {
    energy: [u8; SIZE * SIZE],
}

impl Grid {
    fn parse(input: &str) -> Result<Self> {
        let rows: Vec<&str> = input.lines().collect();
        ensure!(rows.len() == SIZE, "grid must have {SIZE} rows");
        let mut energy = [0u8; SIZE * SIZE];
        for (row, line) in rows.iter().enumerate() {
            ensure!(line.len() == SIZE, "grid row must have {SIZE} columns");
            for (column, ch) in line.chars().enumerate() {
                let digit = ch
                    .to_digit(10)
                    .with_context(|| format!("invalid energy level '{ch}'"))?;
                energy[row * SIZE + column] = u8::try_from(digit).context("digit out of range")?;
            }
        }
        Ok(Self { energy })
    }

    fn neighbors(index: usize) -> impl Iterator<Item = usize> {
        let row = index / SIZE;
        let column = index % SIZE;
        (-1..=1).flat_map(move |dr: isize| {
            (-1..=1).filter_map(move |dc: isize| {
                if dr == 0 && dc == 0 {
                    return None;
                }
                let new_row = row as isize + dr;
                let new_column = column as isize + dc;
                if (0..SIZE as isize).contains(&new_row) && (0..SIZE as isize).contains(&new_column)
                {
                    Some(new_row as usize * SIZE + new_column as usize)
                } else {
                    None
                }
            })
        })
    }

    /// Advances one step, returning the number of octopuses that flashed.
    fn step(&mut self) -> usize {
        let mut flashed = [false; SIZE * SIZE];
        let mut stack: Vec<usize> = Vec::new();

        for index in 0..self.energy.len() {
            self.energy[index] += 1;
            if self.energy[index] > 9 {
                stack.push(index);
            }
        }

        while let Some(index) = stack.pop() {
            if flashed[index] {
                continue;
            }
            flashed[index] = true;
            for neighbor in Self::neighbors(index) {
                self.energy[neighbor] += 1;
                if self.energy[neighbor] > 9 && !flashed[neighbor] {
                    stack.push(neighbor);
                }
            }
        }

        for (energy, flashed) in self.energy.iter_mut().zip(&flashed) {
            if *flashed {
                *energy = 0;
            }
        }

        flashed.iter().filter(|&&f| f).count()
    }
}

/// Returns the total flashes after `steps` steps and the first step where every
/// octopus flashes simultaneously.
fn solve_with(input: &str, steps: usize) -> Result<(usize, usize)> {
    let mut grid = Grid::parse(input)?;
    let total_cells = grid.energy.len();
    let mut total_flashes = 0;
    let mut synchronized_step = None;
    let mut step = 0;
    while synchronized_step.is_none() || step < steps {
        step += 1;
        let flashed = grid.step();
        if step <= steps {
            total_flashes += flashed;
        }
        if flashed == total_cells && synchronized_step.is_none() {
            synchronized_step = Some(step);
        }
    }
    let synchronized_step = synchronized_step.context("no synchronized flash occurred")?;
    Ok((total_flashes, synchronized_step))
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    solve_with(input, 100)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day11-example.txt");
        assert_eq!(solve(input).unwrap(), (1656, 195));
    }

    #[test]
    fn official_example_after_ten_steps() {
        let input = include_str!("../../input/day11-example.txt");
        let (flashes, _) = solve_with(input, 10).unwrap();
        assert_eq!(flashes, 204);
    }

    #[test]
    fn malformed_grid_is_rejected() {
        assert!(solve("not a grid").is_err());
    }

    #[test]
    fn wrong_row_count_is_rejected() {
        let input = "123\n456";
        assert!(solve(input).is_err());
    }
}
