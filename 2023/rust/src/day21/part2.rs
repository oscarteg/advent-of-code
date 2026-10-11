use crate::day21::{custom_error::AocError, garden::Garden};

/// Step budget the elf has on the infinitely tiled map.
const STEPS: u64 = 26_501_365;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let garden = Garden::parse(input)?;
    Ok(garden.reachable_extrapolated(STEPS)?.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    const EXAMPLE: &str = include_str!("../../examples/day21.txt");

    /// Square map with an open border and an open row and column through the
    /// start, the structure the real input has: its tile-aligned counts are
    /// exactly quadratic.
    const STRUCTURED: &str = "\
...........
.#....#....
..#....#...
...#....#..
....#....#.
.....S.....
.#....#....
..#....#...
...#....#..
....#....#.
...........
";

    const OPEN: &str = ".....\n.....\n..S..\n.....\n.....\n";

    #[rstest]
    #[case(6, 16)]
    #[case(10, 50)]
    #[case(50, 1594)]
    #[case(100, 6536)]
    #[case(500, 167004)]
    fn example_step_counts_on_the_tiled_map(#[case] steps: u64, #[case] expected: u64) {
        let garden = Garden::parse(EXAMPLE).expect("example parses");
        assert_eq!(expected, garden.reachable_tiled(steps));
    }

    #[rstest]
    #[case(38, 1309)]
    #[case(49, 2156)]
    #[case(60, 3213)]
    fn extrapolation_matches_the_flood(#[case] steps: u64, #[case] expected: u64) {
        let garden = Garden::parse(STRUCTURED).expect("structured map parses");
        assert_eq!(expected, garden.reachable_tiled(steps));
        assert_eq!(expected, garden.reachable_extrapolated(steps).unwrap());
    }

    #[test]
    fn step_counts_below_two_tiles_are_flooded() -> miette::Result<()> {
        let garden = Garden::parse(EXAMPLE)?;
        assert_eq!(16, garden.reachable_extrapolated(6)?);
        Ok(())
    }

    #[test]
    fn process_fills_the_whole_diamond_without_rocks() -> miette::Result<()> {
        // Without rocks the elf reaches every plot of the diamond with the
        // right parity, which is (steps + 1).pow(2) plots.
        assert_eq!("702322399865956", process(OPEN)?);
        Ok(())
    }

    #[test]
    fn extrapolation_needs_a_square_map() {
        let garden = Garden::parse(".....\n..S..").expect("map parses");
        let error = garden
            .reachable_extrapolated(STEPS)
            .expect_err("a 5x2 map cannot be tiled into whole copies");
        assert!(error.to_string().contains("square"), "{error}");
    }

    #[rstest]
    #[case("")]
    #[case("...\n..")]
    #[case("...\n...\n...")]
    fn malformed_input_is_an_error(#[case] input: &str) {
        assert!(process(input).is_err());
    }
}
