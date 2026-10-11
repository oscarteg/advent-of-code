use crate::day09::{custom_error::AocError, sequence};

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let histories = sequence::parse(input)?;
    let sum: i64 = histories
        .iter()
        .map(|history| history.extrapolate_previous())
        .sum();
    Ok(sum.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day09.txt");
        assert_eq!("2", process(input)?);
        Ok(())
    }

    #[test]
    fn single_history_extrapolates_backward() {
        assert_eq!(process("0 3 6 9 12 15").unwrap(), "-3");
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "1 2 x"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
