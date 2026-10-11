use anyhow::{bail, ensure, Context, Result};
use std::collections::{HashSet, VecDeque};

#[derive(Clone, Copy)]
enum Direction {
    Up,
    Down,
    Left,
    Right,
}

struct Blizzard {
    row: usize,
    col: usize,
    dir: Direction,
}

struct Grid {
    width: usize,
    height: usize,
    start_col: usize,
    end_col: usize,
    blizzards: Vec<Blizzard>,
}

impl Grid {
    fn parse(input: &str) -> Result<Self> {
        let lines: Vec<&str> = input.lines().filter(|line| !line.is_empty()).collect();
        ensure!(
            lines.len() >= 3,
            "grid needs a top wall, an interior and a bottom wall"
        );
        let height = lines.len() - 2;
        let width = lines[0]
            .len()
            .checked_sub(2)
            .filter(|&w| w > 0)
            .context("grid interior must be non-empty")?;

        let top = lines[0];
        let bottom = lines[lines.len() - 1];
        ensure!(top.len() == width + 2, "top wall has the wrong width");
        ensure!(bottom.len() == width + 2, "bottom wall has the wrong width");

        let start_idx = top
            .chars()
            .position(|c| c == '.')
            .context("missing start opening")?;
        let start_col = start_idx
            .checked_sub(1)
            .filter(|&c| c < width)
            .context("start opening must sit above the interior")?;
        let end_idx = bottom
            .chars()
            .position(|c| c == '.')
            .context("missing end opening")?;
        let end_col = end_idx
            .checked_sub(1)
            .filter(|&c| c < width)
            .context("end opening must sit below the interior")?;

        let mut blizzards = Vec::new();
        for (row, line) in lines[1..=height].iter().enumerate() {
            ensure!(line.len() == width + 2, "inconsistent row width");
            for (col_idx, ch) in line.chars().enumerate() {
                if col_idx == 0 || col_idx == width + 1 {
                    ensure!(ch == '#', "side walls must be solid");
                    continue;
                }
                let col = col_idx - 1;
                let dir = match ch {
                    '^' => Some(Direction::Up),
                    'v' => Some(Direction::Down),
                    '<' => Some(Direction::Left),
                    '>' => Some(Direction::Right),
                    '.' => None,
                    other => bail!("unexpected character '{other}' in grid"),
                };
                if let Some(dir) = dir {
                    blizzards.push(Blizzard { row, col, dir });
                }
            }
        }

        Ok(Self {
            width,
            height,
            start_col,
            end_col,
            blizzards,
        })
    }

    /// The set of interior cells occupied by a blizzard at the given time.
    fn occupied_at(&self, time: usize) -> HashSet<(usize, usize)> {
        self.blizzards
            .iter()
            .map(|b| {
                let (dr, dc): (i64, i64) = match b.dir {
                    Direction::Up => (-1, 0),
                    Direction::Down => (1, 0),
                    Direction::Left => (0, -1),
                    Direction::Right => (0, 1),
                };
                let row = (b.row as i64 + dr * time as i64).rem_euclid(self.height as i64) as usize;
                let col = (b.col as i64 + dc * time as i64).rem_euclid(self.width as i64) as usize;
                (row, col)
            })
            .collect()
    }

    /// A position is valid when it is the start opening, the end opening, or
    /// an interior cell.
    fn is_valid(&self, row: i64, col: i64) -> bool {
        if row == -1 {
            return col == self.start_col as i64;
        }
        if row == self.height as i64 {
            return col == self.end_col as i64;
        }
        row >= 0 && col >= 0 && col < self.width as i64
    }
}

fn gcd(a: usize, b: usize) -> usize {
    if b == 0 {
        a
    } else {
        gcd(b, a % b)
    }
}

fn lcm(a: usize, b: usize) -> usize {
    a / gcd(a, b) * b
}

/// BFS over (position, time mod period) states, since the blizzard layout
/// repeats with that period. Returns the elapsed time of the shortest path.
fn shortest_time(
    grid: &Grid,
    occupied: &[HashSet<(usize, usize)>],
    period: usize,
    start_time: usize,
    from: (i64, i64),
    to: (i64, i64),
) -> Result<usize> {
    let mut queue = VecDeque::new();
    let mut visited = HashSet::new();
    queue.push_back((from, start_time));
    visited.insert((from, start_time % period));

    while let Some((pos, time)) = queue.pop_front() {
        if pos == to {
            return Ok(time);
        }
        let next_time = time + 1;
        let occupied_next = &occupied[next_time % period];
        let (row, col) = pos;
        let moves = [(0, 0), (-1, 0), (1, 0), (0, -1), (0, 1)];
        for (dr, dc) in moves {
            let next = (row + dr, col + dc);
            if !grid.is_valid(next.0, next.1) {
                continue;
            }
            if next.0 >= 0
                && next.0 < grid.height as i64
                && occupied_next.contains(&(next.0 as usize, next.1 as usize))
            {
                continue;
            }
            let state = (next, next_time % period);
            if visited.insert(state) {
                queue.push_back((next, next_time));
            }
        }
    }

    bail!("no path found from start to end")
}

/// Parses the grid and returns the elapsed time for each of the three legs
/// (there, back, there again) as part1 and part1 + leg2 + leg3 as part2.
pub fn solve(input: &str) -> Result<(usize, usize)> {
    let grid = Grid::parse(input)?;
    let period = lcm(grid.width, grid.height);
    let occupied: Vec<HashSet<(usize, usize)>> = (0..period).map(|t| grid.occupied_at(t)).collect();

    let start = (-1i64, grid.start_col as i64);
    let end = (grid.height as i64, grid.end_col as i64);

    let there = shortest_time(&grid, &occupied, period, 0, start, end)?;
    let back = shortest_time(&grid, &occupied, period, there, end, start)?;
    let there_again = shortest_time(&grid, &occupied, period, back, start, end)?;

    Ok((there, there_again))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day24-example.txt");
        assert_eq!(solve(input).unwrap(), (18, 54));
    }

    #[test]
    fn no_blizzards_is_a_direct_walk() {
        let input = "#.#\n#.#\n#.#\n#.#\n#.#\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 4);
        assert!(part2 > part1);
    }

    #[test]
    fn missing_openings_is_an_error() {
        let input = "###\n#.#\n#.#\n#.#\n###\n";
        assert!(solve(input).is_err());
    }

    #[test]
    fn inconsistent_row_width_is_an_error() {
        let input = "#.#\n#..#\n#.#\n#.#\n#.#\n";
        assert!(solve(input).is_err());
    }

    #[test]
    fn unknown_character_is_an_error() {
        let input = "#.#\n#x#\n#.#\n#.#\n#.#\n";
        assert!(solve(input).is_err());
    }
}
