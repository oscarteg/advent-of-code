use anyhow::{ensure, Context, Result};
use std::collections::HashSet;

/// Returns the letters answered "yes" by at least one person, and the
/// letters answered "yes" by every person, for one group.
fn group_answers(group: &str) -> Result<(HashSet<char>, HashSet<char>)> {
    let mut people = group.lines();
    let first = people
        .next()
        .context("a group must have at least one person")?;
    ensure!(
        !first.is_empty() && first.chars().all(|c| c.is_ascii_lowercase()),
        "answers must be lowercase letters"
    );
    let mut any: HashSet<char> = first.chars().collect();
    let mut all: HashSet<char> = first.chars().collect();
    for person in people {
        ensure!(
            !person.is_empty() && person.chars().all(|c| c.is_ascii_lowercase()),
            "answers must be lowercase letters"
        );
        let answers: HashSet<char> = person.chars().collect();
        any.extend(&answers);
        all.retain(|c| answers.contains(c));
    }
    Ok((any, all))
}

/// Returns the sum of "anyone answered yes" counts and the sum of
/// "everyone answered yes" counts across all groups.
pub fn solve(input: &str) -> Result<(usize, usize)> {
    let normalized = input.replace("\r\n", "\n");
    let groups: Vec<&str> = normalized.trim().split("\n\n").collect();
    ensure!(!groups.is_empty(), "missing groups");
    let mut any_total = 0;
    let mut all_total = 0;
    for group in groups {
        let (any, all) = group_answers(group)?;
        any_total += any.len();
        all_total += all.len();
    }
    Ok((any_total, all_total))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/2020/day6-example.txt");
        assert_eq!(solve(input).unwrap(), (11, 6));
    }

    #[test]
    fn single_person_group() {
        let input = "abc";
        assert_eq!(solve(input).unwrap(), (3, 3));
    }

    #[test]
    fn windows_newlines_are_accepted() {
        let input = "abc\n\na\nb\nc".replace('\n', "\r\n");
        assert_eq!(solve(&input).unwrap(), (6, 3));
    }

    #[test]
    fn malformed_answers_are_rejected() {
        let input = "abc\n\na1b\nc";
        assert!(solve(input).is_err());
    }

    #[test]
    fn blank_group_is_rejected() {
        let input = "abc\n\n\ndef";
        assert!(solve(input).is_err());
    }
}
