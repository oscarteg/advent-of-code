use crate::day04::custom_error::AocError;
use std::collections::HashSet;

#[derive(Debug)]
pub(crate) struct Card {
    winning_numbers: HashSet<u32>,
    your_numbers: Vec<u32>,
}

impl Card {
    /// Count of this card's numbers that are also winning numbers.
    pub(crate) fn num_matches(&self) -> usize {
        self.your_numbers
            .iter()
            .filter(|n| self.winning_numbers.contains(n))
            .count()
    }

    fn parse(line: &str) -> Result<Self, AocError> {
        let (_, numbers) = line
            .split_once(':')
            .ok_or(AocError::InvalidInput("card is missing ':'"))?;
        let (winning_numbers, your_numbers) = numbers
            .split_once('|')
            .ok_or(AocError::InvalidInput("card is missing '|'"))?;
        let winning_numbers = winning_numbers
            .split_whitespace()
            .map(str::parse)
            .collect::<Result<HashSet<u32>, _>>()?;
        let your_numbers = your_numbers
            .split_whitespace()
            .map(str::parse)
            .collect::<Result<Vec<u32>, _>>()?;
        Ok(Self {
            winning_numbers,
            your_numbers,
        })
    }
}

pub(crate) fn parse(input: &str) -> Result<Vec<Card>, AocError> {
    let cards = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(Card::parse)
        .collect::<Result<Vec<_>, _>>()?;
    if cards.is_empty() {
        return Err(AocError::InvalidInput("no cards"));
    }
    Ok(cards)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn malformed_cards_return_errors() {
        for input in [
            "",
            "Card 1: 1 2 3",
            "Card 1 1 2 | 3 4",
            "Card 1: x 2 | 3 4",
            "Card 1: 1 2 | x 4",
        ] {
            assert!(parse(input).is_err(), "{input}");
        }
    }
}
