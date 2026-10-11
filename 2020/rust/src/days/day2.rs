use anyhow::{bail, Context, Result};
use lazy_static::lazy_static;
use regex::Regex;

lazy_static! {
    // Matches lines like "1-3 a: abcde" into (low, high, letter, password).
    static ref RE: Regex = Regex::new(r"^(\d+)-(\d+) (\w): (\w+)$").unwrap();
}

struct Policy {
    low: usize,
    high: usize,
    letter: char,
    password: String,
}

impl Policy {
    fn parse(line: &str) -> Result<Self> {
        let captures = RE
            .captures(line)
            .with_context(|| format!("line does not match policy format: {line:?}"))?;

        Ok(Policy {
            low: captures[1].parse().context("low bound is not a number")?,
            high: captures[2].parse().context("high bound is not a number")?,
            letter: captures[3]
                .chars()
                .next()
                .context("policy letter is missing")?,
            password: captures[4].to_string(),
        })
    }

    fn is_valid_by_count(&self) -> bool {
        let occurrences = self.password.chars().filter(|&c| c == self.letter).count();
        (self.low..=self.high).contains(&occurrences)
    }

    // Exactly one of the two 1-indexed positions must hold the letter.
    fn is_valid_by_position(&self) -> Result<bool> {
        let at = |pos: usize| -> Result<char> {
            self.password.chars().nth(pos - 1).with_context(|| {
                format!(
                    "password {:?} has no character at position {pos}",
                    self.password
                )
            })
        };

        Ok((at(self.low)? == self.letter) ^ (at(self.high)? == self.letter))
    }
}

fn parse_policies(input: &str) -> Result<Vec<Policy>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(Policy::parse)
        .collect()
}

fn solve_part_01(policies: &[Policy]) -> usize {
    policies.iter().filter(|p| p.is_valid_by_count()).count()
}

fn solve_part_02(policies: &[Policy]) -> Result<usize> {
    let mut valid = 0;
    for policy in policies {
        if policy.is_valid_by_position()? {
            valid += 1;
        }
    }
    Ok(valid)
}

pub fn solve(input: &str) -> Result<(usize, usize)> {
    if input.trim().is_empty() {
        bail!("input is empty");
    }

    let policies = parse_policies(input)?;
    let part1 = solve_part_01(&policies);
    let part2 = solve_part_02(&policies)?;

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day2-example.txt");

    #[test]
    fn sample_part_01() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 2);
    }

    #[test]
    fn sample_part_02() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 1);
    }

    #[test]
    fn position_rule_rejects_both_positions_matching() {
        let policies = parse_policies("2-9 c: ccccccccc").unwrap();
        assert!(!policies[0].is_valid_by_position().unwrap());
    }

    #[test]
    fn malformed_line_is_an_error() {
        assert!(Policy::parse("not a policy line").is_err());
        assert!(solve("not a policy line").is_err());
    }

    #[test]
    fn empty_input_is_an_error() {
        assert!(solve("").is_err());
    }

    #[test]
    fn real_input_answers() {
        let input = std::fs::read_to_string("input/2020/day2.txt").unwrap();
        let (part1, part2) = solve(&input).unwrap();
        assert_eq!(part1, 515);
        assert_eq!(part2, 711);
    }
}
