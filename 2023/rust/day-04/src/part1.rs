use crate::custom_error::AocError;
use miette::Result;
use std::collections::HashSet;

#[tracing::instrument]
pub fn process(input: &str) -> Result<String, AocError> {
    let result = input
        .lines()
        .map(process_line)
        .collect::<Result<Vec<_>, _>>()?
        .into_iter()
        .sum::<u32>()
        .to_string();

    Ok(result)
}

#[tracing::instrument]
pub fn process_line(input: &str) -> Result<u32, AocError> {
    let (_, numbers) = input.split_once(':').unwrap();

    let (winning_numbers, your_numbers) = numbers.split_once('|').unwrap();

    let winning_numbers: HashSet<u32> = winning_numbers
        .split_whitespace()
        .map(|n| n.parse::<u32>().unwrap())
        .collect::<HashSet<u32>>();

    let your_numbers = your_numbers
        .split_whitespace()
        .map(|n| n.parse::<u32>())
        .collect::<Result<Vec<_>, _>>()?
        .into_iter()
        .fold((0, 1), |(acc_points, multiplier), n| {
            if winning_numbers.contains(&n) {
                (multiplier, multiplier * 2)
            } else {
                (acc_points, multiplier)
            }
        })
        .0;

    println!("{your_numbers}");

    return Ok(your_numbers);
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case("Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53", 8)]
    #[case("Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19", 2)]
    #[case("Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1", 2)]
    #[case("Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83", 1)]
    #[case("Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36", 0)]
    #[case("Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11", 0)]
    fn line_test(#[case] line: &str, #[case] expected: u32) {
        assert_eq!(expected, process_line(line).unwrap())
    }

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11";
        assert_eq!("13", process(input)?);
        Ok(())
    }
}
