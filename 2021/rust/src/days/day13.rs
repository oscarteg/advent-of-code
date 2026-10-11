use anyhow::{ensure, Context, Result};
use std::collections::HashSet;

#[derive(Clone, Copy)]
enum Fold {
    X(i32),
    Y(i32),
}

type Dots = HashSet<(i32, i32)>;

fn parse(input: &str) -> Result<(Dots, Vec<Fold>)> {
    let normalized = input.replace("\r\n", "\n");
    let mut sections = normalized.trim().split("\n\n");
    let dots_section = sections.next().context("missing dots")?;
    let folds_section = sections.next().context("missing folds")?;

    let mut dots = HashSet::new();
    for line in dots_section.lines() {
        let (x, y) = line.split_once(',').context("malformed dot")?;
        dots.insert((x.trim().parse()?, y.trim().parse()?));
    }
    ensure!(!dots.is_empty(), "no dots given");

    let mut folds = Vec::new();
    for line in folds_section.lines() {
        let instruction = line
            .strip_prefix("fold along ")
            .context("malformed fold instruction")?;
        let (axis, value) = instruction.split_once('=').context("malformed fold")?;
        let value = value.trim().parse()?;
        folds.push(match axis {
            "x" => Fold::X(value),
            "y" => Fold::Y(value),
            other => anyhow::bail!("unknown fold axis: {other}"),
        });
    }
    ensure!(!folds.is_empty(), "no fold instructions given");

    Ok((dots, folds))
}

fn apply_fold(dots: &Dots, fold: Fold) -> Dots {
    dots.iter()
        .map(|&(x, y)| match fold {
            Fold::X(line) if x > line => (2 * line - x, y),
            Fold::Y(line) if y > line => (x, 2 * line - y),
            _ => (x, y),
        })
        .collect()
}

fn render(dots: &Dots) -> String {
    let width = dots.iter().map(|&(x, _)| x).max().unwrap_or(0) + 1;
    let height = dots.iter().map(|&(_, y)| y).max().unwrap_or(0) + 1;
    (0..height)
        .map(|y| {
            (0..width)
                .map(|x| if dots.contains(&(x, y)) { '#' } else { '.' })
                .collect::<String>()
        })
        .collect::<Vec<_>>()
        .join("\n")
}

/// Returns the dot count after the first fold and the rendered code after every fold.
pub fn solve(input: &str) -> Result<(usize, String)> {
    let (mut dots, folds) = parse(input)?;
    let mut folds = folds.into_iter();
    let first_fold = folds.next().context("no fold instructions given")?;
    dots = apply_fold(&dots, first_fold);
    let part1 = dots.len();
    for fold in folds {
        dots = apply_fold(&dots, fold);
    }
    Ok((part1, render(&dots)))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day13-example.txt");
        let (part1, part2) = solve(input).unwrap();
        assert_eq!(part1, 17);
        assert_eq!(part2, "#####\n#...#\n#...#\n#...#\n#####");
    }

    #[test]
    fn a_single_fold_is_enough_to_compute_part_two() {
        let input = "0,0\n3,0\n\nfold along x=2";
        let (_, part2) = solve(input).unwrap();
        assert_eq!(part2, "##");
    }

    #[test]
    fn overlapping_dots_after_a_fold_count_once() {
        let input = "0,0\n2,0\n1,0\n\nfold along x=1";
        let (part1, _) = solve(input).unwrap();
        assert_eq!(part1, 2);
    }

    #[test]
    fn missing_fold_instructions_are_rejected() {
        assert!(solve("0,0\n1,1").is_err());
    }

    #[test]
    fn malformed_dots_are_rejected() {
        assert!(solve("not-a-dot\n\nfold along x=1").is_err());
    }

    #[test]
    fn unknown_fold_axes_are_rejected() {
        assert!(solve("0,0\n\nfold along z=1").is_err());
    }
}
