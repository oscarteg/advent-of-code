use anyhow::{bail, Context, Result};
use std::str::FromStr;

pub enum Command {
    Forward(u32),
    Down(u32),
    Up(u32),
}

impl FromStr for Command {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self> {
        let (action, value) = s
            .split_once(' ')
            .with_context(|| format!("wrong format: `{s}`"))?;
        let value = value
            .parse::<u32>()
            .with_context(|| format!("invalid distance in `{s}`"))?;
        match action {
            "forward" => Ok(Command::Forward(value)),
            "down" => Ok(Command::Down(value)),
            "up" => Ok(Command::Up(value)),
            _ => bail!("unknown direction: `{s}`"),
        }
    }
}

fn parse(input: &str) -> Result<Vec<Command>> {
    input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(Command::from_str)
        .collect()
}

/// Part 1: forward moves horizontally, up/down change depth directly.
fn part1(commands: &[Command]) -> u32 {
    let (horizontal, depth) = commands.iter().fold((0u32, 0u32), |(h, d), cmd| match cmd {
        Command::Forward(v) => (h + v, d),
        Command::Down(v) => (h, d + v),
        Command::Up(v) => (h, d.saturating_sub(*v)),
    });
    horizontal * depth
}

/// Part 2: up/down change aim; forward moves horizontally and adjusts depth by aim * distance.
fn part2(commands: &[Command]) -> u32 {
    let (horizontal, depth, _) = commands
        .iter()
        .fold((0u32, 0u32, 0u32), |(h, d, aim), cmd| match cmd {
            Command::Forward(v) => (h + v, d + aim * v, aim),
            Command::Down(v) => (h, d, aim + v),
            Command::Up(v) => (h, d, aim.saturating_sub(*v)),
        });
    horizontal * depth
}

pub fn solve(input: &str) -> Result<(u32, u32)> {
    let commands = parse(input)?;
    Ok((part1(&commands), part2(&commands)))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day2-example.txt");

    #[test]
    fn part1_example() {
        let commands = parse(EXAMPLE).unwrap();
        assert_eq!(part1(&commands), 150);
    }

    #[test]
    fn part2_example() {
        let commands = parse(EXAMPLE).unwrap();
        assert_eq!(part2(&commands), 900);
    }

    #[test]
    fn part1_real_input() {
        let input = std::fs::read_to_string("input/day2.txt").unwrap();
        let commands = parse(&input).unwrap();
        assert_eq!(part1(&commands), 2039256);
    }

    #[test]
    fn part2_real_input() {
        let input = std::fs::read_to_string("input/day2.txt").unwrap();
        let commands = parse(&input).unwrap();
        assert_eq!(part2(&commands), 1856459736);
    }

    #[test]
    fn up_saturates_at_zero() {
        let commands = parse("up 5").unwrap();
        assert_eq!(part1(&commands), 0);
        assert_eq!(part2(&commands), 0);
    }

    #[test]
    fn malformed_input_returns_err() {
        assert!(parse("sideways 5").is_err());
        assert!(parse("forward").is_err());
        assert!(parse("forward abc").is_err());
    }
}
