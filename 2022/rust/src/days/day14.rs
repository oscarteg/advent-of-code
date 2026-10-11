use anyhow::{bail, ensure, Context, Result};
use std::collections::HashSet;

type Point = (i64, i64);

const SOURCE: Point = (500, 0);

fn parse_point(text: &str) -> Result<Point> {
    let (x, y) = text
        .trim()
        .split_once(',')
        .with_context(|| format!("point `{text}` is missing a comma"))?;
    Ok((
        x.trim()
            .parse()
            .with_context(|| format!("invalid x in `{text}`"))?,
        y.trim()
            .parse()
            .with_context(|| format!("invalid y in `{text}`"))?,
    ))
}

fn parse_rocks(input: &str) -> Result<HashSet<Point>> {
    let mut rocks = HashSet::new();
    for line in input.lines() {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        let points = line
            .split("->")
            .map(parse_point)
            .collect::<Result<Vec<Point>>>()?;
        ensure!(points.len() >= 2, "a rock path needs at least two points");
        for pair in points.windows(2) {
            let (x1, y1) = pair[0];
            let (x2, y2) = pair[1];
            if x1 == x2 {
                for y in y1.min(y2)..=y1.max(y2) {
                    rocks.insert((x1, y));
                }
            } else if y1 == y2 {
                for x in x1.min(x2)..=x1.max(x2) {
                    rocks.insert((x, y1));
                }
            } else {
                bail!("rock path segment from ({x1},{y1}) to ({x2},{y2}) is not axis-aligned");
            }
        }
    }
    ensure!(!rocks.is_empty(), "no rock paths found");
    Ok(rocks)
}

/// Drops sand from the source until it either falls past `max_y` (no floor) or
/// comes to rest on the source itself (with a floor), returning the count of
/// settled units.
fn simulate(rocks: &HashSet<Point>, max_y: i64, floor_y: Option<i64>) -> u64 {
    let mut settled: HashSet<Point> = rocks.clone();
    let mut count = 0u64;
    while !settled.contains(&SOURCE) {
        let mut sand = SOURCE;
        loop {
            if floor_y.is_none() && sand.1 > max_y {
                return count;
            }
            let below = (sand.0, sand.1 + 1);
            let below_left = (sand.0 - 1, sand.1 + 1);
            let below_right = (sand.0 + 1, sand.1 + 1);
            let is_blocked = |point: Point| settled.contains(&point) || floor_y == Some(point.1);
            if !is_blocked(below) {
                sand = below;
            } else if !is_blocked(below_left) {
                sand = below_left;
            } else if !is_blocked(below_right) {
                sand = below_right;
            } else {
                break;
            }
        }
        settled.insert(sand);
        count += 1;
    }
    count
}

pub fn solve(input: &str) -> Result<(u64, u64)> {
    let rocks = parse_rocks(input)?;
    let max_y = rocks
        .iter()
        .map(|&(_, y)| y)
        .max()
        .context("no rocks to bound the cave")?;
    let part1 = simulate(&rocks, max_y, None);
    let part2 = simulate(&rocks, max_y, Some(max_y + 2));
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day14-example.txt");
        assert_eq!(solve(input).unwrap(), (24, 93));
    }

    #[test]
    fn single_horizontal_rock_catches_a_bounded_pile() {
        // A flat ledge below the source catches sand in a triangular pile
        // until the pile grows past the ledge's edges and spills into the abyss.
        let input = "498,3 -> 502,3";
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 4);
    }

    #[test]
    fn missing_comma_is_rejected() {
        assert!(parse_rocks("498 4 -> 498 6").is_err());
    }

    #[test]
    fn diagonal_segment_is_rejected() {
        assert!(parse_rocks("498,4 -> 500,6").is_err());
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("").is_err());
    }
}
