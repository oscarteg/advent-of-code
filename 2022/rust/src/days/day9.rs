use anyhow::{bail, ensure, Context, Result};
use std::collections::HashSet;

#[derive(Debug, Clone, Copy)]
struct Motion {
    dx: i32,
    dy: i32,
    steps: u32,
}

fn parse(input: &str) -> Result<Vec<Motion>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(|line| {
            let (direction, steps) = line
                .split_once(' ')
                .with_context(|| format!("malformed motion line: {line:?}"))?;
            let steps = steps
                .parse()
                .with_context(|| format!("invalid step count in line: {line:?}"))?;
            let (dx, dy) = match direction {
                "R" => (1, 0),
                "L" => (-1, 0),
                "U" => (0, 1),
                "D" => (0, -1),
                other => bail!("unknown direction: {other:?}"),
            };
            Ok(Motion { dx, dy, steps })
        })
        .collect()
}

fn follow(head: (i32, i32), tail: (i32, i32)) -> (i32, i32) {
    let (dx, dy) = (head.0 - tail.0, head.1 - tail.1);
    if dx.abs() <= 1 && dy.abs() <= 1 {
        return tail;
    }
    (tail.0 + dx.signum(), tail.1 + dy.signum())
}

fn simulate(motions: &[Motion], rope_len: usize) -> usize {
    let mut knots = vec![(0i32, 0i32); rope_len];
    let mut visited: HashSet<(i32, i32)> = HashSet::new();
    visited.insert(knots[rope_len - 1]);

    for motion in motions {
        for _ in 0..motion.steps {
            knots[0].0 += motion.dx;
            knots[0].1 += motion.dy;
            for i in 1..rope_len {
                knots[i] = follow(knots[i - 1], knots[i]);
            }
            visited.insert(knots[rope_len - 1]);
        }
    }

    visited.len()
}

/// Returns the number of positions the tail of a 2-knot rope and a 10-knot
/// rope each visit at least once.
pub fn solve(input: &str) -> Result<(usize, usize)> {
    let motions = parse(input)?;
    ensure!(!motions.is_empty(), "input contains no motions");
    Ok((simulate(&motions, 2), simulate(&motions, 10)))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example_part1() {
        let input = include_str!("../../input/day9-example.txt");
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 13);
    }

    #[test]
    fn official_example_part2() {
        let input = include_str!("../../input/day9-example.txt");
        let (_, part2) = solve(input).unwrap();
        assert_eq!(part2, 1);
    }

    #[test]
    fn larger_example_part2() {
        let input = "R 5\nU 8\nL 8\nD 3\nR 17\nD 10\nL 25\nU 20\n";
        let (_, part2) = solve(input).unwrap();
        assert_eq!(part2, 36);
    }

    #[test]
    fn single_step_does_not_move_tail() {
        let (part1, part2) = solve("R 1\n").unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 1);
    }

    #[test]
    fn empty_input_is_rejected() {
        assert!(solve("").is_err());
    }

    #[test]
    fn malformed_direction_is_rejected() {
        assert!(solve("X 4\n").is_err());
    }

    #[test]
    fn malformed_step_count_is_rejected() {
        assert!(solve("R four\n").is_err());
    }
}
