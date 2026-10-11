use anyhow::{anyhow, ensure};

pub fn solve(input: &str) -> anyhow::Result<usize> {
    let mut grid = parse_grid(input)?;
    let height = grid.len();
    let width = if height > 0 { grid[0].len() } else { 0 };
    ensure!(width > 0, "Grid must have at least one column");

    let mut step = 0;
    loop {
        step += 1;
        let moved = simulate_step(&mut grid, height, width);
        if !moved {
            return Ok(step);
        }
    }
}

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum Cell {
    Empty,
    East,
    South,
}

fn parse_grid(input: &str) -> anyhow::Result<Vec<Vec<Cell>>> {
    input
        .lines()
        .filter(|line| !line.is_empty())
        .map(|line| {
            line.chars()
                .map(|ch| match ch {
                    '.' => Ok(Cell::Empty),
                    '>' => Ok(Cell::East),
                    'v' => Ok(Cell::South),
                    c => Err(anyhow!("Invalid character: {}", c)),
                })
                .collect()
        })
        .collect()
}

fn simulate_step(grid: &mut Vec<Vec<Cell>>, height: usize, width: usize) -> bool {
    let moved_east = move_east(grid, height, width);
    let moved_south = move_south(grid, height, width);
    moved_east || moved_south
}

fn move_east(grid: &mut Vec<Vec<Cell>>, height: usize, width: usize) -> bool {
    let mut moved = false;
    let mut new_grid = grid.clone();
    for row in 0..height {
        for col in 0..width {
            if grid[row][col] == Cell::East {
                let next_col = (col + 1) % width;
                if grid[row][next_col] == Cell::Empty {
                    new_grid[row][col] = Cell::Empty;
                    new_grid[row][next_col] = Cell::East;
                    moved = true;
                }
            }
        }
    }
    *grid = new_grid;
    moved
}

fn move_south(grid: &mut Vec<Vec<Cell>>, height: usize, width: usize) -> bool {
    let mut moved = false;
    let mut new_grid = grid.clone();
    for row in 0..height {
        for col in 0..width {
            if grid[row][col] == Cell::South {
                let next_row = (row + 1) % height;
                if grid[next_row][col] == Cell::Empty {
                    new_grid[row][col] = Cell::Empty;
                    new_grid[next_row][col] = Cell::South;
                    moved = true;
                }
            }
        }
    }
    *grid = new_grid;
    moved
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_example() {
        let input = include_str!("../../input/day25-example.txt");
        assert_eq!(solve(input).unwrap(), 58);
    }

    #[test]
    fn packed_grid_stops_on_first_step() {
        assert_eq!(solve(">v\nv>").unwrap(), 1);
    }

    #[test]
    fn east_herd_moves_only_into_free_cells() {
        let mut grid = parse_grid("...>>>>>...").unwrap();
        move_east(&mut grid, 1, 11);
        assert_eq!(grid, parse_grid("...>>>>.>..").unwrap());
    }

    #[test]
    fn test_invalid_char() {
        let input = "x";
        assert!(solve(input).is_err());
    }

    #[test]
    fn test_empty_grid() {
        let input = "";
        assert!(solve(input).is_err());
    }
}
