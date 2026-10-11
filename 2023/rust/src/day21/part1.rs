use crate::day21::{custom_error::AocError, garden::Garden};

/// Step budget the elf has on the map as handed out.
const STEPS: u64 = 64;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let garden = Garden::parse(input)?;
    Ok(garden.reachable(STEPS).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    const EXAMPLE: &str = include_str!("../../examples/day21.txt");

    #[rstest]
    #[case(0, 1)]
    #[case(1, 2)]
    #[case(2, 4)]
    #[case(3, 6)]
    #[case(6, 16)]
    fn example_step_counts(#[case] steps: u64, #[case] expected: u64) {
        let garden = Garden::parse(EXAMPLE).expect("example parses");
        assert_eq!(expected, garden.reachable(steps));
    }

    #[test]
    fn process_saturates_the_small_example() -> miette::Result<()> {
        // 64 steps outgrow the 11x11 example, leaving every even-distance plot.
        assert_eq!("42", process(EXAMPLE)?);
        Ok(())
    }

    #[rstest]
    #[case(0, 1)]
    #[case(1, 0)]
    #[case(6, 0)]
    fn a_walled_in_start_goes_nowhere(#[case] steps: u64, #[case] expected: u64) {
        let garden = Garden::parse("###\n#S#\n###").expect("map parses");
        assert_eq!(expected, garden.reachable(steps));
    }

    #[rstest]
    #[case("")]
    #[case("...\n..")]
    #[case("...\n.S.\n.X.")]
    #[case("...\n...\n...")]
    #[case("S..\n...\n..S")]
    fn malformed_input_is_an_error(#[case] input: &str) {
        assert!(process(input).is_err());
    }
}
