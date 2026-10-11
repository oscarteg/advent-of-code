use crate::day04::{custom_error::AocError, scratchcards};

/// Sums each card's score: 0 matches score 0, otherwise `2^(matches - 1)`.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let cards = scratchcards::parse(input)?;
    let total: u32 = cards
        .iter()
        .map(|card| match card.num_matches() {
            0 => 0,
            matches => 2u32.pow(matches as u32 - 1),
        })
        .sum();
    Ok(total.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case("Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53", "8")]
    #[case("Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19", "2")]
    #[case("Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36", "0")]
    fn single_card_score(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day04.txt");
        assert_eq!("13", process(input)?);
        Ok(())
    }

    #[test]
    fn no_matches_scores_zero() {
        let input = "Card 1: 1 2 3 | 4 5 6";
        assert_eq!(process(input).unwrap(), "0");
    }

    #[test]
    fn malformed_input_returns_error() {
        assert!(process("Card 1: 1 2 3").is_err());
    }
}
