use crate::day18::custom_error::AocError;
use crate::day18::dig_plan::{lagoon_area, DigPlan};

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let plan = DigPlan::parse(input, false)?;
    Ok(lagoon_area(&plan.steps).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day18.txt");
        assert_eq!("62", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("R 4 (#000000)\nD 4 (#000000)\nL 4 (#000000)\nU 4 (#000000)", "25")]
    #[case("R 1 (#000000)\nD 1 (#000000)\nL 1 (#000000)\nU 1 (#000000)", "4")]
    fn small_loops(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            "X 4 (#000000)",
            "R x (#000000)",
            "R 4",
            "RR 4 (#000000)",
        ] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
