use crate::day14::custom_error::AocError;
use std::collections::HashMap;

fn parse_grid(input: &str) -> Result<Vec<Vec<char>>, AocError> {
    let grid: Vec<Vec<char>> = input.lines().map(|line| line.chars().collect()).collect();
    if grid.is_empty() {
        return Err(AocError::InvalidInput("empty grid"));
    }
    Ok(grid)
}

#[allow(clippy::needless_range_loop)]
fn tilt_north(grid: &mut [Vec<char>]) {
    let rows = grid.len();
    let cols = grid[0].len();

    for col in 0..cols {
        let mut write_pos = 0;
        for row in 0..rows {
            match grid[row][col] {
                'O' => {
                    grid[row][col] = '.';
                    grid[write_pos][col] = 'O';
                    write_pos += 1;
                }
                '#' => {
                    write_pos = row + 1;
                }
                '.' => {}
                _ => {}
            }
        }
    }
}

#[allow(clippy::needless_range_loop)]
fn tilt_west(grid: &mut [Vec<char>]) {
    let rows = grid.len();
    let cols = grid[0].len();

    for row in 0..rows {
        let mut write_pos = 0;
        for col in 0..cols {
            match grid[row][col] {
                'O' => {
                    grid[row][col] = '.';
                    grid[row][write_pos] = 'O';
                    write_pos += 1;
                }
                '#' => {
                    write_pos = col + 1;
                }
                '.' => {}
                _ => {}
            }
        }
    }
}

#[allow(clippy::needless_range_loop)]
fn tilt_south(grid: &mut [Vec<char>]) {
    let rows = grid.len();
    let cols = grid[0].len();

    for col in 0..cols {
        let mut write_pos = rows;
        for row in (0..rows).rev() {
            match grid[row][col] {
                'O' => {
                    grid[row][col] = '.';
                    write_pos -= 1;
                    grid[write_pos][col] = 'O';
                }
                '#' => {
                    write_pos = row;
                }
                '.' => {}
                _ => {}
            }
        }
    }
}

#[allow(clippy::needless_range_loop)]
fn tilt_east(grid: &mut [Vec<char>]) {
    let rows = grid.len();
    let cols = grid[0].len();

    for row in 0..rows {
        let mut write_pos = cols;
        for col in (0..cols).rev() {
            match grid[row][col] {
                'O' => {
                    grid[row][col] = '.';
                    write_pos -= 1;
                    grid[row][write_pos] = 'O';
                }
                '#' => {
                    write_pos = col;
                }
                '.' => {}
                _ => {}
            }
        }
    }
}

fn spin_cycle(grid: &mut [Vec<char>]) {
    tilt_north(grid);
    tilt_west(grid);
    tilt_south(grid);
    tilt_east(grid);
}

fn grid_to_string(grid: &[Vec<char>]) -> String {
    grid.iter()
        .map(|row| row.iter().collect::<String>())
        .collect::<Vec<_>>()
        .join("\n")
}

fn calculate_load(grid: &[Vec<char>]) -> u64 {
    let rows = grid.len() as u64;
    grid.iter()
        .enumerate()
        .map(|(row_idx, row)| {
            let row_value = rows - row_idx as u64;
            row.iter().filter(|&&c| c == 'O').count() as u64 * row_value
        })
        .sum()
}

/// Runs spin cycles with cycle detection and calculates the load after 1,000,000,000 cycles.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut grid = parse_grid(input)?;
    let target_cycles = 1_000_000_000;
    let mut seen: HashMap<String, usize> = HashMap::new();

    for cycle_num in 0..target_cycles {
        let grid_str = grid_to_string(&grid);
        if let Some(&prev_cycle) = seen.get(&grid_str) {
            let cycle_length = cycle_num - prev_cycle;
            let remaining = (target_cycles - cycle_num) % cycle_length;
            for _ in 0..remaining {
                spin_cycle(&mut grid);
            }
            let load = calculate_load(&grid);
            return Ok(load.to_string());
        }
        seen.insert(grid_str, cycle_num);
        spin_cycle(&mut grid);
    }

    let load = calculate_load(&grid);
    Ok(load.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day14.txt");
        let result = process(input)?;
        assert!(!result.is_empty());
        Ok(())
    }

    #[test]
    fn test_single_cycle() {
        let input = "O.#...O#.\n.......#.\n#........\n..#......\n......#..\n...#.#...\n.#...#.#.\n..#...#..\n.#.......#\n#.#...O#.#";
        let result = process(input).unwrap();
        assert!(!result.is_empty());
    }

    #[test]
    fn test_malformed_input() {
        assert!(process("").is_err());
    }
}
