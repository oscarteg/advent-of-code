use crate::day07::custom_error::AocError;
use crate::day07::hands::Hand;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let mut hands: Vec<Hand> = input
        .lines()
        .filter(|line| !line.is_empty())
        .map(Hand::parse)
        .collect::<Result<_, _>>()?;

    hands.sort();

    let total: u32 = hands
        .iter()
        .enumerate()
        .map(|(rank, hand)| hand.bid * (rank as u32 + 1))
        .sum();

    Ok(total.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day07.txt");
        assert_eq!("6440", process(input)?);
        Ok(())
    }

    #[test]
    fn test_single_hand() {
        let input = "32T3K 765";
        assert_eq!("765", process(input).unwrap());
    }

    #[test]
    fn test_invalid_hand() {
        let input = "INVALID 100";
        assert!(process(input).is_err());
    }

    #[test]
    fn test_empty_input() -> miette::Result<()> {
        let input = "";
        assert_eq!("0", process(input)?);
        Ok(())
    }
}
