use anyhow::{ensure, Result};
use std::collections::HashMap;

#[derive(Debug, Clone, Copy)]
struct Line {
    x1: i32,
    y1: i32,
    x2: i32,
    y2: i32,
}

impl Line {
    fn parse(s: &str) -> Result<Self> {
        let parts: Vec<&str> = s.split("->").collect();
        ensure!(parts.len() == 2, "line must contain exactly one '->'");

        let (p1, p2) = (parts[0].trim(), parts[1].trim());

        let coords1: Vec<i32> = p1
            .split(',')
            .map(|c| c.trim().parse())
            .collect::<Result<Vec<_>, _>>()?;
        ensure!(coords1.len() == 2, "first coordinate must have x,y");

        let coords2: Vec<i32> = p2
            .split(',')
            .map(|c| c.trim().parse())
            .collect::<Result<Vec<_>, _>>()?;
        ensure!(coords2.len() == 2, "second coordinate must have x,y");

        Ok(Self {
            x1: coords1[0],
            y1: coords1[1],
            x2: coords2[0],
            y2: coords2[1],
        })
    }

    fn is_horizontal_or_vertical(self) -> bool {
        self.x1 == self.x2 || self.y1 == self.y2
    }

    fn is_diagonal(self) -> bool {
        (self.x1 - self.x2).abs() == (self.y1 - self.y2).abs()
    }

    fn points_on_line(self) -> Vec<(i32, i32)> {
        let dx = (self.x2 - self.x1).signum();
        let dy = (self.y2 - self.y1).signum();
        let steps = (self.x2 - self.x1).abs().max((self.y2 - self.y1).abs());

        (0..=steps)
            .map(|i| (self.x1 + dx * i, self.y1 + dy * i))
            .collect()
    }
}

/// Returns (part1, part2) - counts of overlapping points.
/// Part 1 considers only horizontal and vertical lines.
/// Part 2 includes diagonal lines at 45-degree angles.
pub fn solve(input: &str) -> Result<(usize, usize)> {
    let lines: Vec<Line> = input
        .lines()
        .filter(|l| !l.trim().is_empty())
        .map(Line::parse)
        .collect::<Result<Vec<_>>>()?;

    let mut grid_hv: HashMap<(i32, i32), usize> = HashMap::new();
    let mut grid_all: HashMap<(i32, i32), usize> = HashMap::new();

    for line in &lines {
        for point in line.points_on_line() {
            if line.is_horizontal_or_vertical() {
                *grid_hv.entry(point).or_insert(0) += 1;
            }
            if line.is_horizontal_or_vertical() || line.is_diagonal() {
                *grid_all.entry(point).or_insert(0) += 1;
            }
        }
    }

    let part1 = grid_hv.values().filter(|&&count| count >= 2).count();
    let part2 = grid_all.values().filter(|&&count| count >= 2).count();

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day5-example.txt");
        assert_eq!(solve(input).unwrap(), (5, 12));
    }

    #[test]
    fn horizontal_line_overlap() {
        let input = "0,0 -> 2,0\n1,0 -> 3,0\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 2);
        assert_eq!(part2, 2);
    }

    #[test]
    fn vertical_line_overlap() {
        let input = "0,0 -> 0,3\n0,1 -> 0,2\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 2);
        assert_eq!(part2, 2);
    }

    #[test]
    fn diagonal_lines_part2_only() {
        let input = "0,0 -> 2,2\n2,0 -> 0,2\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 1);
    }

    #[test]
    fn single_point_overlap() {
        let input = "0,0 -> 0,0\n0,0 -> 0,0\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 1);
        assert_eq!(part2, 1);
    }

    #[test]
    fn non_axis_aligned_non_diagonal_ignored() {
        let input = "0,0 -> 2,1\n0,1 -> 2,0\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 0);
    }

    #[test]
    fn malformed_input_missing_arrow() {
        let input = "0,0 5,5";
        assert!(solve(input).is_err());
    }

    #[test]
    fn malformed_input_invalid_coordinate() {
        let input = "0,x -> 5,5";
        assert!(solve(input).is_err());
    }

    #[test]
    fn empty_input() {
        let input = "";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 0);
        assert_eq!(part2, 0);
    }

    #[test]
    fn empty_lines_ignored() {
        let input = "0,0 -> 0,1\n\n0,0 -> 0,1\n";
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 2);
        assert_eq!(part2, 2);
    }
}
