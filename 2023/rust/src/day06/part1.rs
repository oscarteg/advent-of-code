use crate::day06::{
    custom_error::AocError,
    races::{parse_races, Race},
};

/// Multiplies the number of winning hold times across every race.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let races = parse_races(input)?;
    let product: u64 = races.iter().map(Race::winning_hold_times).product();
    Ok(product.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(Race { time: 7, distance: 9 }, 4)]
    #[case(Race { time: 15, distance: 40 }, 8)]
    #[case(Race { time: 30, distance: 200 }, 9)]
    fn winning_hold_times_matches_example(#[case] race: Race, #[case] expected: u64) {
        assert_eq!(expected, race.winning_hold_times());
    }

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day06.txt");
        assert_eq!("288", process(input)?);
        Ok(())
    }

    #[test]
    fn no_hold_time_wins_when_record_is_unbeatable() {
        let input = "Time: 1\nDistance: 100";
        assert_eq!(process(input).unwrap(), "0");
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "Time: 1\nDistance: 1 2", "Time: a\nDistance: 1"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
