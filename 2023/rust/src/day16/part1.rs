use crate::day16::custom_error::AocError;
use std::collections::HashSet;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
enum Direction {
    Up,
    Down,
    Left,
    Right,
}

impl Direction {
    fn dx(&self) -> isize {
        match self {
            Direction::Up => -1,
            Direction::Down => 1,
            Direction::Left => 0,
            Direction::Right => 0,
        }
    }

    fn dy(&self) -> isize {
        match self {
            Direction::Up => 0,
            Direction::Down => 0,
            Direction::Left => -1,
            Direction::Right => 1,
        }
    }
}

fn reflect_forward(dir: Direction) -> Direction {
    match dir {
        Direction::Up => Direction::Right,
        Direction::Down => Direction::Left,
        Direction::Left => Direction::Down,
        Direction::Right => Direction::Up,
    }
}

fn reflect_backslash(dir: Direction) -> Direction {
    match dir {
        Direction::Up => Direction::Left,
        Direction::Down => Direction::Right,
        Direction::Left => Direction::Up,
        Direction::Right => Direction::Down,
    }
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let grid: Vec<Vec<char>> = input.lines().map(|line| line.chars().collect()).collect();

    if grid.is_empty() || grid[0].is_empty() {
        return Ok("0".to_string());
    }

    let energized = simulate(&grid, 0, 0, Direction::Right);
    Ok(energized.len().to_string())
}

fn simulate(
    grid: &[Vec<char>],
    start_row: usize,
    start_col: usize,
    start_dir: Direction,
) -> HashSet<(usize, usize)> {
    let mut visited = HashSet::new();
    let mut queue = vec![(start_row as isize, start_col as isize, start_dir)];
    let mut seen = HashSet::new();

    while let Some((row, col, dir)) = queue.pop() {
        if row < 0 || col < 0 || row >= grid.len() as isize || col >= grid[0].len() as isize {
            continue;
        }

        let row_usize = row as usize;
        let col_usize = col as usize;

        if !seen.insert((row_usize, col_usize, dir)) {
            continue;
        }

        visited.insert((row_usize, col_usize));

        let cell = grid[row_usize][col_usize];
        match cell {
            '.' => {
                let next_row = row + dir.dx();
                let next_col = col + dir.dy();
                queue.push((next_row, next_col, dir));
            }
            '/' => {
                let new_dir = reflect_forward(dir);
                let next_row = row + new_dir.dx();
                let next_col = col + new_dir.dy();
                queue.push((next_row, next_col, new_dir));
            }
            '\\' => {
                let new_dir = reflect_backslash(dir);
                let next_row = row + new_dir.dx();
                let next_col = col + new_dir.dy();
                queue.push((next_row, next_col, new_dir));
            }
            '-' => match dir {
                Direction::Left | Direction::Right => {
                    let next_row = row + dir.dx();
                    let next_col = col + dir.dy();
                    queue.push((next_row, next_col, dir));
                }
                Direction::Up | Direction::Down => {
                    queue.push((row, col - 1, Direction::Left));
                    queue.push((row, col + 1, Direction::Right));
                }
            },
            '|' => match dir {
                Direction::Up | Direction::Down => {
                    let next_row = row + dir.dx();
                    let next_col = col + dir.dy();
                    queue.push((next_row, next_col, dir));
                }
                Direction::Left | Direction::Right => {
                    queue.push((row - 1, col, Direction::Up));
                    queue.push((row + 1, col, Direction::Down));
                }
            },
            _ => {}
        }
    }

    visited
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day16.txt");
        assert_eq!("46", process(input)?);
        Ok(())
    }

    #[test]
    fn test_empty_grid() -> miette::Result<()> {
        let input = "";
        let result = process(input)?;
        assert_eq!("0", result);
        Ok(())
    }

    #[test]
    fn test_single_dot() -> miette::Result<()> {
        let input = ".";
        let result = process(input)?;
        assert_eq!("1", result);
        Ok(())
    }
}
