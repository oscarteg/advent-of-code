use crate::day12::{custom_error::AocError, springs::parse_records};

/// Every row of the real records is folded into five copies of itself.
const COPIES: usize = 5;

/// Sums the arrangements of each condition record unfolded five times.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let arrangements: u64 = parse_records(input)?
        .iter()
        .map(|record| record.unfold(COPIES).arrangements())
        .sum();
    Ok(arrangements.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case("???.### 1,1,3", "1")]
    #[case(".??..??...?##. 1,1,3", "16384")]
    #[case("?#?#?#?#?#?#?#? 1,3,1,6", "1")]
    #[case("????.#...#... 4,1,1", "16")]
    #[case("????.######..#####. 1,6,5", "2500")]
    #[case("?###???????? 3,2,1", "506250")]
    #[case("# 1", "1")]
    #[case("? 1", "1")]
    #[case("#.# 2", "0")]
    fn unfolded_row_arrangements(#[case] row: &str, #[case] expected: &str) -> miette::Result<()> {
        assert_eq!(expected, process(row)?);
        Ok(())
    }

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day12.txt");
        assert_eq!("525152", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("")]
    #[case("???.###")]
    #[case("??!.### 1,1,3")]
    #[case("???.### 1,x,3")]
    #[case("???.### 1,0,3")]
    fn malformed_records_are_errors(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
