use std::{cmp, collections::HashMap, str::FromStr};

use crate::custom_error::AocError;

#[derive(Debug, Eq, Hash, PartialEq)]
enum Color {
    Red,
    Green,
    Blue,
}

impl FromStr for Color {
    type Err = AocError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "red" => Ok(Color::Red),
            "green" => Ok(Color::Green),
            "blue" => Ok(Color::Blue),
            _ => unreachable!(),
        }
    }
}

#[tracing::instrument]
fn process_line(line: &str) -> u32 {
    let parts = line.split(": ").nth(1).unwrap().split("; ");
    let mut min_cubes: HashMap<Color, u32> = HashMap::new();

    for part in parts {
        for p in part.split(", ") {
            let x: Vec<&str> = p.split_whitespace().collect();

            let number = x[0].parse::<u32>().unwrap();
            let color = x[1].parse::<Color>().unwrap();

            let current_min = min_cubes.entry(color).or_insert(number);

            *current_min = cmp::max(*current_min, number);
        }
    }

    min_cubes.values().product::<u32>()
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let result = input.lines().fold(0, |acc, line| {
        let result = process_line(line);

        acc + result
    });

    Ok(result.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green";
        assert_eq!("2286", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green", 48)]
    fn line_test(#[case] line: &str, #[case] expected: u32) {
        assert_eq!(expected, process_line(line))
    }
}
