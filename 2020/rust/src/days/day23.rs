use anyhow::{ensure, Context, Result};
use std::iter;

const PART1_MOVES: usize = 100;
const PART2_CUPS: usize = 1_000_000;
const PART2_MOVES: usize = 10_000_000;
/// With fewer cups the three picked-up cups would overlap the current cup,
/// which the rules do not define.
const MIN_CUPS: usize = 5;

/// Parses the single line of cup labels into a permutation of `1..=n`, the
/// invariant the successor array relies on.
fn parse(input: &str) -> Result<Vec<usize>> {
    let mut lines = input.lines().map(str::trim).filter(|line| !line.is_empty());
    let labels = lines.next().context("expected a line of cup labels")?;
    ensure!(
        lines.next().is_none(),
        "expected a single line of cup labels"
    );

    let count = labels.chars().count();
    let mut cups = Vec::with_capacity(count);
    let mut seen = vec![false; count + 1];
    for character in labels.chars() {
        let digit = character
            .to_digit(10)
            .with_context(|| format!("invalid cup label: {character:?}"))?;
        let label = usize::try_from(digit)?;
        ensure!(
            (1..=count).contains(&label) && !seen[label],
            "cup labels must be a permutation of 1..={count}, got {labels:?}"
        );
        seen[label] = true;
        cups.push(label);
    }
    ensure!(
        cups.len() >= MIN_CUPS,
        "need at least {MIN_CUPS} cups, got {}",
        cups.len()
    );
    Ok(cups)
}

/// Plays `moves` moves on a circle of `total` cups: the labels of `cups`
/// clockwise, followed by `cups.len() + 1..=total`. Returns the successor
/// array, in which `next[label]` is the label immediately clockwise of it.
fn play(cups: &[usize], total: usize, moves: usize) -> Result<Vec<usize>> {
    let &start = cups.first().context("circle without cups")?;
    ensure!(
        total >= cups.len(),
        "circle of {total} cups cannot hold {} labels",
        cups.len()
    );

    // Every label starts out followed by the next one, which already wires the
    // padding cups; the input labels and the closing link are patched in below.
    let mut next: Vec<usize> = (0..=total).map(|label| label + 1).collect();
    let tail = if total > cups.len() {
        cups.len() + 1
    } else {
        start
    };
    for (&cup, &clockwise) in cups
        .iter()
        .zip(cups.iter().skip(1).chain(iter::once(&tail)))
    {
        next[cup] = clockwise;
    }
    if total > cups.len() {
        next[total] = start;
    }

    let mut current = start;
    for _ in 0..moves {
        let first = next[current];
        let second = next[first];
        let third = next[second];
        next[current] = next[third];

        let mut destination = counter_clockwise_label(current, total);
        while destination == first || destination == second || destination == third {
            destination = counter_clockwise_label(destination, total);
        }

        next[third] = next[destination];
        next[destination] = first;
        current = next[current];
    }
    Ok(next)
}

fn counter_clockwise_label(label: usize, total: usize) -> usize {
    if label == 1 {
        total
    } else {
        label - 1
    }
}

fn labels_after_one(next: &[usize]) -> String {
    let mut labels = String::new();
    let mut cup = next[1];
    while cup != 1 {
        labels.push_str(&cup.to_string());
        cup = next[cup];
    }
    labels
}

fn product_after_one(next: &[usize]) -> Result<u64> {
    let first = next[1];
    let second = next[first];
    Ok(u64::try_from(first)? * u64::try_from(second)?)
}

pub fn solve(input: &str) -> Result<(String, u64)> {
    let cups = parse(input)?;
    let part1 = labels_after_one(&play(&cups, cups.len(), PART1_MOVES)?);
    let part2 = product_after_one(&play(&cups, PART2_CUPS, PART2_MOVES)?)?;
    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day23-example.txt");

    fn example_cups() -> Vec<usize> {
        parse(EXAMPLE).expect("example parses")
    }

    #[test]
    fn part1_example() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, "67384529");
    }

    #[test]
    fn part2_example() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 149_245_887_792);
    }

    #[test]
    fn ten_moves_match_the_walkthrough() {
        let cups = example_cups();
        let next = play(&cups, cups.len(), 10).unwrap();
        assert_eq!(labels_after_one(&next), "92658374");
    }

    #[test]
    fn one_move_splices_three_cups_behind_the_destination() {
        let cups = example_cups();
        let next = play(&cups, cups.len(), 1).unwrap();
        // After the first move the circle reads 3 2 8 9 1 5 4 6 7.
        assert_eq!(labels_after_one(&next), "54673289");
    }

    #[test]
    fn zero_moves_keep_a_small_circle_in_order() {
        let cups = parse("32415").unwrap();
        let next = play(&cups, cups.len(), 0).unwrap();
        assert_eq!(labels_after_one(&next), "5324");
    }

    #[test]
    fn padding_cups_continue_the_circle() {
        let cups = example_cups();
        let next = play(&cups, 12, 0).unwrap();
        let mut circle = Vec::new();
        let mut cup = cups[0];
        for _ in 0..12 {
            circle.push(cup);
            cup = next[cup];
        }
        assert_eq!(circle, [3, 8, 9, 1, 2, 5, 4, 6, 7, 10, 11, 12]);
        assert_eq!(cup, cups[0], "the circle closes onto the first cup");
    }

    #[test]
    fn malformed_input_is_rejected() {
        for input in [
            "",
            "   ",
            "38912546x",
            "389125460",
            "388125467",
            "98765",
            "3241",
            "389125467\n213456789",
        ] {
            assert!(parse(input).is_err(), "expected {input:?} to be rejected");
        }
    }
}
