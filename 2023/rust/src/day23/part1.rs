use crate::day23::{
    custom_error::AocError,
    trails::{Map, Slopes},
};

/// The longest hike when every slope has to be left downhill.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let steps = Map::parse(input)?.longest_hike(Slopes::Icy)?;
    Ok(steps.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    const CORRIDOR: &str = "#.#\n#.#\n#.#\n";
    const DOWNHILL_SLOPE: &str = "#.#\n#v#\n#.#\n";
    const UPHILL_SLOPE: &str = "#.#\n#^#\n#.#\n";
    const DEAD_END_SPUR: &str = "#.##\n#..#\n#.##\n";
    // Two ways around one loop: three steps along the left side, five along the
    // right, whose last tile is a slope pointing back west.
    const LOOP_WITH_SLOPE: &str = "#.###\n#.<.#\n#.#.#\n#...#\n##.##\n";

    #[rstest]
    #[case::official_example(include_str!("../../examples/day23.txt"), "94")]
    #[case::straight_corridor(CORRIDOR, "2")]
    #[case::downhill_slope(DOWNHILL_SLOPE, "2")]
    #[case::dead_end_spur(DEAD_END_SPUR, "2")]
    #[case::slope_blocks_the_long_way_round(LOOP_WITH_SLOPE, "5")]
    #[case::entrance_is_the_exit("#.#\n", "0")]
    fn longest_hike(#[case] input: &str, #[case] expected: &str) -> miette::Result<()> {
        assert_eq!(expected, process(input)?);
        Ok(())
    }

    #[test]
    fn uphill_slope_blocks_the_only_route() {
        assert!(process(UPHILL_SLOPE).is_err());
    }

    #[rstest]
    #[case::empty("")]
    #[case::ragged_rows("#.#\n#.##\n#.#\n")]
    #[case::unknown_tile("#.#\n#x#\n#.#\n")]
    #[case::no_entrance("###\n#.#\n#.#\n")]
    #[case::no_exit("#.#\n#.#\n###\n")]
    fn rejects_malformed_input(#[case] input: &str) {
        assert!(process(input).is_err());
    }
}
