use crate::day02::{cube_game::parse, custom_error::AocError};

#[tracing::instrument]
pub fn process(input: &str, red: u32, green: u32, blue: u32) -> miette::Result<String, AocError> {
    let games = parse(input)?;
    let sum: u32 = games
        .iter()
        .filter(|game| game.is_possible(red, green, blue))
        .map(|game| game.id)
        .sum();
    Ok(sum.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn official_example() -> miette::Result<()> {
        let input = include_str!("../../examples/day02.txt");
        assert_eq!("8", process(input, 12, 13, 14)?);
        Ok(())
    }

    #[rstest]
    #[case("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green", true)]
    #[case("Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green", false)]
    fn game_possibility(#[case] line: &str, #[case] expected: bool) {
        let games = parse(line).unwrap();
        assert_eq!(games[0].is_possible(12, 13, 14), expected);
    }

    #[test]
    fn exact_limits_are_possible() -> miette::Result<()> {
        let input = "Game 1: 12 red, 13 green, 14 blue";
        assert_eq!("1", process(input, 12, 13, 14)?);
        Ok(())
    }

    #[test]
    fn malformed_input_is_an_error() {
        assert!(process("", 12, 13, 14).is_err());
        assert!(process("Game 1: 3 purple", 12, 13, 14).is_err());
    }
}
