use crate::day09::{custom_error::AocError, sequence};

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let histories = sequence::parse(input)?;
    let sum: i64 = histories
        .iter()
        .map(|history| history.extrapolate_next())
        .sum();
    Ok(sum.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day09.txt");
        assert_eq!("114", process(input)?);
        Ok(())
    }

    #[test]
    fn single_history_sums_to_its_own_extrapolation() {
        assert_eq!(process("1 1 1").unwrap(), "1");
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "1 2 x"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
