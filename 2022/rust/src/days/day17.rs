use anyhow::{bail, ensure, Result};
use std::collections::HashMap;

const CHAMBER_WIDTH: i32 = 7;
const PART1_ROCKS: u64 = 2022;
const PART2_ROCKS: u64 = 1_000_000_000_000;
const CYCLE_PROFILE_DEPTH: usize = 50;

/// Rock shapes as bottom-to-top row bitmasks (bit 0 is the leftmost column), with their width.
const ROCK_SHAPES: [(&[u8], i32); 5] = [
    (&[0b1111], 4),
    (&[0b010, 0b111, 0b010], 3),
    (&[0b111, 0b100, 0b100], 3),
    (&[0b1, 0b1, 0b1, 0b1], 1),
    (&[0b11, 0b11], 2),
];

fn parse_jets(input: &str) -> Result<Vec<i32>> {
    let trimmed = input.trim();
    ensure!(!trimmed.is_empty(), "jet pattern must not be empty");
    trimmed
        .chars()
        .map(|c| match c {
            '<' => Ok(-1),
            '>' => Ok(1),
            other => bail!("invalid jet character: {other}"),
        })
        .collect()
}

fn collides(shape: &[u8], width: i32, x: i32, y: i32, rows: &[u8]) -> bool {
    if x < 0 || x + width > CHAMBER_WIDTH || y < 0 {
        return true;
    }
    shape.iter().enumerate().any(|(i, &row_mask)| {
        let row_y = y + i as i32;
        row_y < rows.len() as i32 && (rows[row_y as usize] as u16) & ((row_mask as u16) << x) != 0
    })
}

fn settle(shape: &[u8], x: i32, y: i32, rows: &mut Vec<u8>) {
    for (i, &row_mask) in shape.iter().enumerate() {
        let row_y = (y + i as i32) as usize;
        while rows.len() <= row_y {
            rows.push(0);
        }
        rows[row_y] |= row_mask << x;
    }
}

/// Top `CYCLE_PROFILE_DEPTH` rows of the tower, used as part of the cycle-detection key.
fn surface_profile(rows: &[u8]) -> Vec<u8> {
    let depth = CYCLE_PROFILE_DEPTH.min(rows.len());
    rows[rows.len() - depth..].to_vec()
}

/// Height of the tower after dropping `target` rocks, using cycle detection on
/// (rock shape index, jet index, surface profile) so huge `target` values stay fast.
fn tower_height(jets: &[i32], target: u64) -> u64 {
    let mut rows: Vec<u8> = Vec::new();
    let mut jet_index = 0usize;
    let mut heights: Vec<u64> = vec![0];
    let mut seen: HashMap<(usize, usize, Vec<u8>), (u64, u64)> = HashMap::new();
    let mut n: u64 = 0;

    while n < target {
        let rock_idx = (n % 5) as usize;
        let (shape, width) = ROCK_SHAPES[rock_idx];
        let mut x: i32 = 2;
        let mut y: i32 = rows.len() as i32 + 3;
        loop {
            let dir = jets[jet_index];
            jet_index = (jet_index + 1) % jets.len();
            let nx = x + dir;
            if !collides(shape, width, nx, y, &rows) {
                x = nx;
            }
            if !collides(shape, width, x, y - 1, &rows) {
                y -= 1;
            } else {
                settle(shape, x, y, &mut rows);
                break;
            }
        }
        n += 1;
        heights.push(rows.len() as u64);

        let key = (rock_idx, jet_index, surface_profile(&rows));
        if let Some(&(prev_n, prev_height)) = seen.get(&key) {
            let cycle_len = n - prev_n;
            let cycle_gain = rows.len() as u64 - prev_height;
            let remaining = target - n;
            let cycles = remaining / cycle_len;
            let remainder = remaining % cycle_len;
            let remainder_gain = heights[(prev_n + remainder) as usize] - prev_height;
            return rows.len() as u64 + cycles * cycle_gain + remainder_gain;
        }
        seen.insert(key, (n, rows.len() as u64));
    }
    rows.len() as u64
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let jets = parse_jets(input)?;
    let part1 = tower_height(&jets, PART1_ROCKS);
    let part2 = tower_height(&jets, PART2_ROCKS);
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day17-example.txt");
        assert_eq!(solve(input).unwrap(), (3068, 1_514_285_714_288));
    }

    #[test]
    fn single_jet_pushes_rock_to_the_wall() {
        let jets = parse_jets("<").unwrap();
        assert_eq!(tower_height(&jets, 1), 1);
    }

    #[test]
    fn malformed_jet_pattern_is_rejected() {
        assert!(solve("abc").is_err());
        assert!(solve("").is_err());
    }
}
