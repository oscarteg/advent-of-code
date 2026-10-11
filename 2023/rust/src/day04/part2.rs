use crate::day04::{custom_error::AocError, scratchcards};

/// Each card's matches win one copy of each of the following `matches` cards.
/// Tracks a running copy count per card and propagates it forward.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let cards = scratchcards::parse(input)?;
    let mut copies = vec![1u64; cards.len()];
    for (index, card) in cards.iter().enumerate() {
        let matches = card.num_matches();
        let count = copies[index];
        for offset in 1..=matches {
            if let Some(target) = copies.get_mut(index + offset) {
                *target += count;
            }
        }
    }
    Ok(copies.iter().sum::<u64>().to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day04.txt");
        assert_eq!("30", process(input)?);
        Ok(())
    }

    #[test]
    fn matches_past_the_last_card_are_ignored() {
        let input = "Card 1: 1 2 | 1 2\nCard 2: 1 2 | 3 4";
        assert_eq!(process(input).unwrap(), "3");
    }

    #[test]
    fn no_matches_keeps_a_single_copy_per_card() {
        let input = "Card 1: 1 2 | 3 4\nCard 2: 1 2 | 3 4";
        assert_eq!(process(input).unwrap(), "2");
    }

    #[test]
    fn malformed_input_returns_error() {
        assert!(process("Card 1: 1 2 | x").is_err());
    }
}
