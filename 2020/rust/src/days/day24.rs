use anyhow::{bail, ensure, Result};
use std::collections::HashMap;

type Pos = (i32, i32);

#[derive(Clone)]
struct Grid {
    black: HashMap<Pos, bool>,
}

impl Grid {
    fn new() -> Self {
        Grid {
            black: HashMap::new(),
        }
    }

    fn is_black(&self, pos: Pos) -> bool {
        self.black.get(&pos).copied().unwrap_or(false)
    }

    fn set_black(&mut self, pos: Pos, black: bool) {
        if !black {
            self.black.remove(&pos);
        } else {
            self.black.insert(pos, true);
        }
    }

    fn count_black(&self) -> usize {
        self.black.values().filter(|&&b| b).count()
    }

    fn neighbors(&self, pos: Pos) -> Vec<Pos> {
        let (x, y) = pos;
        vec![
            (x + 1, y),
            (x - 1, y),
            (x, y + 1),
            (x, y - 1),
            (x + 1, y - 1),
            (x - 1, y + 1),
        ]
    }

    fn count_black_neighbors(&self, pos: Pos) -> usize {
        self.neighbors(pos)
            .iter()
            .filter(|&&p| self.is_black(p))
            .count()
    }

    fn step(&self) -> Grid {
        let mut new_grid = Grid::new();

        let mut to_check = std::collections::HashSet::new();
        for &pos in self.black.keys() {
            to_check.insert(pos);
            for &neighbor in &self.neighbors(pos) {
                to_check.insert(neighbor);
            }
        }

        for pos in to_check {
            let is_black = self.is_black(pos);
            let black_neighbors = self.count_black_neighbors(pos);

            let should_be_black = if is_black {
                black_neighbors == 1 || black_neighbors == 2
            } else {
                black_neighbors == 2
            };

            if should_be_black {
                new_grid.set_black(pos, true);
            }
        }

        new_grid
    }
}

fn parse_path(line: &str) -> Result<Pos> {
    let mut x = 0i32;
    let mut y = 0i32;

    let bytes = line.as_bytes();
    let mut i = 0;

    while i < bytes.len() {
        let ch = bytes[i] as char;
        match ch {
            'e' => {
                x += 1;
                i += 1;
            }
            'w' => {
                x -= 1;
                i += 1;
            }
            'n' => {
                ensure!(i + 1 < bytes.len(), "Incomplete direction after 'n'");
                let next = bytes[i + 1] as char;
                if next == 'e' {
                    y += 1;
                    i += 2;
                } else if next == 'w' {
                    x -= 1;
                    y += 1;
                    i += 2;
                } else {
                    bail!("Invalid direction after 'n': {}", next);
                }
            }
            's' => {
                ensure!(i + 1 < bytes.len(), "Incomplete direction after 's'");
                let next = bytes[i + 1] as char;
                if next == 'e' {
                    x += 1;
                    y -= 1;
                    i += 2;
                } else if next == 'w' {
                    y -= 1;
                    i += 2;
                } else {
                    bail!("Invalid direction after 's': {}", next);
                }
            }
            _ => {
                bail!("Invalid direction character: {}", ch);
            }
        }
    }

    Ok((x, y))
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let mut grid = Grid::new();

    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        let pos = parse_path(line)?;
        let is_black = grid.is_black(pos);
        grid.set_black(pos, !is_black);
    }

    let part1 = grid.count_black();

    for _ in 0..100 {
        grid = grid.step();
    }

    let part2 = grid.count_black();

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_path() {
        assert_eq!(parse_path("e").unwrap(), (1, 0));
        assert_eq!(parse_path("w").unwrap(), (-1, 0));
        assert_eq!(parse_path("ne").unwrap(), (0, 1));
        assert_eq!(parse_path("nw").unwrap(), (-1, 1));
        assert_eq!(parse_path("se").unwrap(), (1, -1));
        assert_eq!(parse_path("sw").unwrap(), (0, -1));
        assert_eq!(parse_path("esew").unwrap(), (1, -1));
        assert_eq!(parse_path("nwwswee").unwrap(), (0, 0));
    }

    #[test]
    fn test_part1() {
        let example = include_str!("../../input/2020/day24-example.txt");
        let (part1, _) = solve(example).unwrap();
        assert_eq!(part1, 10);
    }

    #[test]
    fn test_part2() {
        let example = include_str!("../../input/2020/day24-example.txt");
        let (_, part2) = solve(example).unwrap();
        assert_eq!(part2, 2208);
    }

    #[test]
    fn test_malformed_input() {
        assert!(parse_path("nq").is_err());
        assert!(parse_path("sn").is_err());
        assert!(parse_path("n").is_err());
        assert!(parse_path("s").is_err());
    }

    #[test]
    fn test_empty_input() {
        let (part1, part2) = solve("").unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 0);
    }

    #[test]
    fn test_single_tile() {
        let (part1, _) = solve("e").unwrap();
        assert_eq!(part1, 1);
    }

    #[test]
    fn test_toggle_back() {
        let (part1, _) = solve("e\ne").unwrap();
        assert_eq!(part1, 0);
    }
}
