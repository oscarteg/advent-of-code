use crate::day02::{cube_game::parse, custom_error::AocError};

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let games = parse(input)?;
    let sum: u32 = games.iter().map(|game| game.power()).sum();
    Ok(sum.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn official_example() -> miette::Result<()> {
        let input = include_str!("../../examples/day02.txt");
        assert_eq!("2286", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green", 48)]
    #[case(
        "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
        1560
    )]
    fn game_power(#[case] line: &str, #[case] expected: u32) {
        let games = parse(line).unwrap();
        assert_eq!(games[0].power(), expected);
    }

    #[test]
    fn single_color_game_has_zero_power() -> miette::Result<()> {
        let input = "Game 1: 3 red; 5 red";
        assert_eq!("0", process(input)?);
        Ok(())
    }

    #[test]
    fn malformed_input_is_an_error() {
        assert!(process("").is_err());
        assert!(process("Game 1: three red").is_err());
    }
}
