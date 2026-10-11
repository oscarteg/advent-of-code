use crate::day11::custom_error::AocError;
use crate::day11::part1::solve_with_expansion;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let answer = solve_with_expansion(input, 1_000_000)?;
    Ok(answer.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day11.txt");
        let result = process(input)?;
        assert_eq!("82000210", result);
        Ok(())
    }

    #[test]
    fn test_malformed_input() {
        assert!(process("").is_err());
    }
}
