use crate::day06::{custom_error::AocError, races::parse_single_race};

/// Joins the digits on each line into a single race and counts its winning hold times.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let race = parse_single_race(input)?;
    Ok(race.winning_hold_times().to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day06.txt");
        assert_eq!("71503", process(input)?);
        Ok(())
    }

    #[test]
    fn digits_are_joined_without_separating_spaces() {
        let race = parse_single_race("Time: 1 2\nDistance: 3 4").unwrap();
        assert_eq!(race.time, 12);
        assert_eq!(race.distance, 34);
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "Time: 1", "Time: a\nDistance: 1"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
