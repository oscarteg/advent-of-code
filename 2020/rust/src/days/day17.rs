use anyhow::{bail, Result};
use std::collections::HashSet;

type Pos3 = (i32, i32, i32);
type Pos4 = (i32, i32, i32, i32);

fn active_cells(input: &str) -> Result<Vec<(i32, i32)>> {
    let mut cells = Vec::new();
    for (y, line) in input.lines().enumerate() {
        for (x, ch) in line.chars().enumerate() {
            match ch {
                '#' => cells.push((x as i32, y as i32)),
                '.' => {}
                _ => bail!("line {}: unexpected {ch:?}, want '#' or '.'", y + 1),
            }
        }
    }
    Ok(cells)
}

fn neighbors_3d(pos: Pos3) -> Vec<Pos3> {
    let (x, y, z) = pos;
    let mut result = Vec::new();
    for dx in [-1, 0, 1] {
        for dy in [-1, 0, 1] {
            for dz in [-1, 0, 1] {
                if dx != 0 || dy != 0 || dz != 0 {
                    result.push((x + dx, y + dy, z + dz));
                }
            }
        }
    }
    result
}

fn neighbors_4d(pos: Pos4) -> Vec<Pos4> {
    let (x, y, z, w) = pos;
    let mut result = Vec::new();
    for dx in [-1, 0, 1] {
        for dy in [-1, 0, 1] {
            for dz in [-1, 0, 1] {
                for dw in [-1, 0, 1] {
                    if dx != 0 || dy != 0 || dz != 0 || dw != 0 {
                        result.push((x + dx, y + dy, z + dz, w + dw));
                    }
                }
            }
        }
    }
    result
}

fn count_active_neighbors_3d(pos: Pos3, active: &HashSet<Pos3>) -> usize {
    neighbors_3d(pos)
        .iter()
        .filter(|p| active.contains(p))
        .count()
}

fn count_active_neighbors_4d(pos: Pos4, active: &HashSet<Pos4>) -> usize {
    neighbors_4d(pos)
        .iter()
        .filter(|p| active.contains(p))
        .count()
}

fn step_3d(active: &HashSet<Pos3>) -> HashSet<Pos3> {
    let mut to_check = HashSet::new();
    for &pos in active {
        to_check.insert(pos);
        for neighbor in neighbors_3d(pos) {
            to_check.insert(neighbor);
        }
    }

    let mut new_active = HashSet::new();
    for pos in to_check {
        let count = count_active_neighbors_3d(pos, active);
        let is_active = active.contains(&pos);

        if (is_active && (count == 2 || count == 3)) || (!is_active && count == 3) {
            new_active.insert(pos);
        }
    }
    new_active
}

fn step_4d(active: &HashSet<Pos4>) -> HashSet<Pos4> {
    let mut to_check = HashSet::new();
    for &pos in active {
        to_check.insert(pos);
        for neighbor in neighbors_4d(pos) {
            to_check.insert(neighbor);
        }
    }

    let mut new_active = HashSet::new();
    for pos in to_check {
        let count = count_active_neighbors_4d(pos, active);
        let is_active = active.contains(&pos);

        if (is_active && (count == 2 || count == 3)) || (!is_active && count == 3) {
            new_active.insert(pos);
        }
    }
    new_active
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    let cells = active_cells(input.trim())?;

    let mut active_3d: HashSet<Pos3> = cells.iter().map(|&(x, y)| (x, y, 0)).collect();
    for _ in 0..6 {
        active_3d = step_3d(&active_3d);
    }
    let part1 = active_3d.len();

    let mut active_4d: HashSet<Pos4> = cells.iter().map(|&(x, y)| (x, y, 0, 0)).collect();
    for _ in 0..6 {
        active_4d = step_4d(&active_4d);
    }
    let part2 = active_4d.len();

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_example() {
        let input = include_str!("../../input/2020/day17-example.txt");
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 112);
        assert_eq!(part2, 848);
    }

    #[test]
    fn test_single_cube() {
        let input = "#";
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 0);
    }

    #[test]
    fn test_empty_input() {
        let (part1, part2) = solve("").unwrap();
        assert_eq!((part1, part2), (0, 0));
    }

    #[test]
    fn test_malformed_input() {
        assert!(solve(".#.\n.x.\n###").is_err());
    }
}
