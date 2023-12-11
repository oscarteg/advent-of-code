use crate::custom_error::AocError;
use std::{cmp, collections::HashMap, str::FromStr};

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

            _ => Err(AocError::ParseError {
                input: s.to_string(),
                source: s.to_string(), // if you have source code, replace with it
                span: (0..s.len()).into(),
            }),
        }
    }
}

#[tracing::instrument]
fn process_line(line: &str) -> u32 {
    let parts = line.split(": ").nth(1).unwrap().split("; ");

    let min_cubes = parts.fold(HashMap::new(), |mut acc, part| {
        part.split(", ").for_each(|p| {
            let (number, color) = p.split_once(' ').unwrap();
            let number = number.parse::<u32>().unwrap();
            let color = color.parse::<Color>().unwrap();

            acc.entry(color)
                .and_modify(|e| *e = cmp::max(*e, number))
                .or_insert(number);
        });
        acc
    });

    min_cubes.values().product()
}

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let result = input.lines().map(process_line).sum::<u32>();

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
