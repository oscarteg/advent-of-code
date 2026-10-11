use anyhow::{bail, ensure, Result};
use std::str::FromStr;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Move {
    Rock = 1,
    Paper = 2,
    Scissors = 3,
}

impl Move {
    fn defeats(self, other: Move) -> bool {
        matches!(
            (self, other),
            (Move::Rock, Move::Scissors)
                | (Move::Paper, Move::Rock)
                | (Move::Scissors, Move::Paper)
        )
    }

    fn vs(self, other: Move) -> Outcome {
        if self == other {
            Outcome::Draw
        } else if self.defeats(other) {
            Outcome::Win
        } else {
            Outcome::Loss
        }
    }

    fn for_outcome(opponent: Move, outcome: Outcome) -> Move {
        match outcome {
            Outcome::Loss => match opponent {
                Move::Rock => Move::Scissors,
                Move::Paper => Move::Rock,
                Move::Scissors => Move::Paper,
            },
            Outcome::Draw => opponent,
            Outcome::Win => match opponent {
                Move::Rock => Move::Paper,
                Move::Paper => Move::Scissors,
                Move::Scissors => Move::Rock,
            },
        }
    }

    fn score(self, outcome: Outcome) -> u32 {
        let outcome_score = outcome.score();
        (self as u32) + outcome_score
    }
}

impl FromStr for Move {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self> {
        match s {
            "A" | "X" => Ok(Move::Rock),
            "B" | "Y" => Ok(Move::Paper),
            "C" | "Z" => Ok(Move::Scissors),
            _ => bail!("invalid move: {}", s),
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Outcome {
    Loss,
    Draw,
    Win,
}

impl Outcome {
    fn score(self) -> u32 {
        match self {
            Outcome::Loss => 0,
            Outcome::Draw => 3,
            Outcome::Win => 6,
        }
    }
}

impl FromStr for Outcome {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self> {
        match s {
            "X" => Ok(Outcome::Loss),
            "Y" => Ok(Outcome::Draw),
            "Z" => Ok(Outcome::Win),
            _ => bail!("invalid outcome: {}", s),
        }
    }
}

struct Round {
    opponent: Move,
    second: String,
}

impl Round {
    fn parse(line: &str) -> Result<Self> {
        let parts: Vec<&str> = line.split_whitespace().collect();
        ensure!(parts.len() == 2, "invalid round format");
        Ok(Round {
            opponent: parts[0].parse()?,
            second: parts[1].to_string(),
        })
    }

    fn score_part1(&self) -> Result<u32> {
        let our_move: Move = self.second.parse()?;
        let outcome = our_move.vs(self.opponent);
        Ok(our_move.score(outcome))
    }

    fn score_part2(&self) -> Result<u32> {
        let outcome: Outcome = self.second.parse()?;
        let our_move = Move::for_outcome(self.opponent, outcome);
        Ok(our_move.score(outcome))
    }
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let lines = input
        .trim()
        .lines()
        .filter(|line| !line.is_empty())
        .collect::<Vec<_>>();

    let mut part1 = 0u32;
    let mut part2 = 0u32;

    for line in lines {
        let round = Round::parse(line)?;
        part1 += round.score_part1()?;
        part2 += round.score_part2()?;
    }

    Ok((part1, part2))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "A Y
B X
C Z";

    #[test]
    fn part1_example() {
        let (part1, _) = solve(EXAMPLE).unwrap();
        assert_eq!(part1, 15);
    }

    #[test]
    fn part2_example() {
        let (_, part2) = solve(EXAMPLE).unwrap();
        assert_eq!(part2, 12);
    }

    #[test]
    fn part1_single_round_win() {
        let (part1, _) = solve("A Y").unwrap();
        assert_eq!(part1, 8);
    }

    #[test]
    fn part1_single_round_loss() {
        let (part1, _) = solve("B X").unwrap();
        assert_eq!(part1, 1);
    }

    #[test]
    fn part1_single_round_draw() {
        let (part1, _) = solve("C Z").unwrap();
        assert_eq!(part1, 6);
    }

    #[test]
    fn part2_single_round_loss() {
        let (_, part2) = solve("B X").unwrap();
        assert_eq!(part2, 1);
    }

    #[test]
    fn part2_single_round_draw() {
        let (_, part2) = solve("C Y").unwrap();
        assert_eq!(part2, 6);
    }

    #[test]
    fn part2_single_round_win() {
        let (_, part2) = solve("A Z").unwrap();
        assert_eq!(part2, 8);
    }

    #[test]
    fn malformed_move() {
        let result = solve("D X");
        assert!(result.is_err());
    }

    #[test]
    fn malformed_format() {
        let result = solve("A");
        assert!(result.is_err());
    }

    #[test]
    fn real_input() {
        let input = std::fs::read_to_string("input/2020/day2.txt")
            .or_else(|_| std::fs::read_to_string("input/day2.txt"))
            .unwrap_or_default();
        if !input.is_empty() {
            let result = solve(&input);
            assert!(result.is_ok());
            let (part1, part2) = result.unwrap();
            assert_eq!(part1, 13809);
            assert_eq!(part2, 12316);
        }
    }
}
