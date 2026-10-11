use crate::day15::custom_error::AocError;
use crate::day15::hashmap::{hash, steps};

/// Sums the HASH algorithm result of every step in the initialization sequence.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let sum: u32 = steps(input).map(hash).sum();
    Ok(sum.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day15.txt");
        assert_eq!("1320", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("HASH", "52")]
    #[case("rn=1,cm-", "283")]
    fn process_sums_known_steps(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[test]
    fn empty_input_sums_to_zero() {
        assert_eq!("0", process("").unwrap());
    }
}
