use anyhow::{bail, ensure, Context, Result};
use std::collections::HashSet;

/// Returns the priority of an item: a-z are 1-26, A-Z are 27-52.
fn priority(c: char) -> Result<u32> {
    if c.is_ascii_lowercase() {
        Ok(c as u32 - 'a' as u32 + 1)
    } else if c.is_ascii_uppercase() {
        Ok(c as u32 - 'A' as u32 + 27)
    } else {
        bail!("invalid item type: {c:?}")
    }
}

/// Finds the single item shared by every one of the given strings.
fn common_item(groups: &[&str]) -> Result<char> {
    ensure!(!groups.is_empty(), "no compartments/rucksacks to compare");

    let mut common: HashSet<char> = groups[0].chars().collect();
    for group in &groups[1..] {
        let chars: HashSet<char> = group.chars().collect();
        common.retain(|c| chars.contains(c));
    }

    let mut iter = common.into_iter();
    let item = iter
        .next()
        .context("no item common to all compartments/rucksacks")?;
    ensure!(iter.next().is_none(), "more than one common item found");
    Ok(item)
}

fn part_one(input: &str) -> Result<u32> {
    let mut sum = 0;
    for rucksack in input.lines() {
        ensure!(
            rucksack.len() % 2 == 0 && !rucksack.is_empty(),
            "rucksack must have an even, non-zero number of items: {rucksack:?}"
        );
        let (first, second) = rucksack.split_at(rucksack.len() / 2);
        let item = common_item(&[first, second])
            .with_context(|| format!("rucksack {rucksack:?} has no shared item"))?;
        sum += priority(item)?;
    }
    Ok(sum)
}

fn part_two(input: &str) -> Result<u32> {
    let lines: Vec<&str> = input.lines().collect();
    ensure!(
        lines.len().is_multiple_of(3) && !lines.is_empty(),
        "input must contain a multiple of 3 rucksacks"
    );

    let mut sum = 0;
    for group in lines.chunks(3) {
        let badge = common_item(group)
            .with_context(|| format!("group {group:?} has no shared badge item"))?;
        sum += priority(badge)?;
    }
    Ok(sum)
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    Ok((part_one(input)?, part_two(input)?))
}

#[cfg(test)]
mod tests {
    use super::*;

    const INPUT: &str = include_str!("../../input/day3-example.txt");

    #[test]
    fn test_part_1() -> Result<()> {
        assert_eq!(part_one(INPUT)?, 157);

        let file = std::fs::read_to_string("input/day3.txt")?;
        assert_eq!(part_one(&file)?, 7845);
        Ok(())
    }

    #[test]
    fn test_part_2() -> Result<()> {
        assert_eq!(part_two(INPUT)?, 70);

        let file = std::fs::read_to_string("input/day3.txt")?;
        assert_eq!(part_two(&file)?, 2790);
        Ok(())
    }

    #[test]
    fn test_part_1_rejects_odd_length_rucksack() {
        assert!(part_one("abc").is_err());
    }

    #[test]
    fn test_part_2_rejects_non_multiple_of_three() {
        assert!(part_two("abc\ndef").is_err());
    }

    #[test]
    fn test_priority_rejects_non_letter() {
        assert!(priority('1').is_err());
    }

    #[test]
    fn test_solve_example() -> Result<()> {
        assert_eq!(solve(INPUT)?, (157, 70));
        Ok(())
    }
}
