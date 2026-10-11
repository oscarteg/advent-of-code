use crate::day18::custom_error::AocError;
use crate::day18::dig_plan::{lagoon_area, DigPlan};

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let plan = DigPlan::parse(input, true)?;
    Ok(lagoon_area(&plan.steps).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day18.txt");
        assert_eq!("952408144115", process(input)?);
        Ok(())
    }

    #[test]
    fn small_square() {
        let input = "R 4 (#000040)\nD 4 (#000041)\nL 4 (#000042)\nU 4 (#000043)";
        assert_eq!("25", process(input).unwrap());
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            "R 4 (#000000",
            "R 4 (#00000)",
            "R 4 (#00000g)",
            "R 4 000000",
        ] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
