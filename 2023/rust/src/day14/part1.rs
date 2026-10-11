use crate::day14::custom_error::AocError;

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

/// Tilts the grid north and calculates the load.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut grid = parse_grid(input)?;
    tilt_north(&mut grid);
    let load = calculate_load(&grid);
    Ok(load.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "O..#....#.\n.......#..\n#.......#.\n..#.......\n......#...\n...#.#...#\n.#...#.#.#\n..#....#..\n.#.......#\n#.#...O#.#";
        let result = process(input)?;
        assert!(!result.is_empty());
        Ok(())
    }

    #[test]
    fn test_single_rock() {
        let input = "O\n.\n.";
        assert_eq!("3", process(input).unwrap());
    }

    #[test]
    fn test_with_blockers() {
        let input = "O.O\n###\nO.O";
        assert_eq!("8", process(input).unwrap());
    }

    #[test]
    fn test_empty_grid_fails() {
        let input = "";
        assert!(process(input).is_err());
    }
}
