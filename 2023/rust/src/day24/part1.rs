use crate::day24::{custom_error::AocError, hailstone};

/// The puzzle's test area: both coordinates within these bounds, inclusive.
const TEST_AREA: (i128, i128) = (200_000_000_000_000, 400_000_000_000_000);

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    process_with_area(input, TEST_AREA.0, TEST_AREA.1)
}

/// Counts the crossings inside an arbitrary test area; the example uses `7..=27`.
#[tracing::instrument]
pub fn process_with_area(input: &str, low: i128, high: i128) -> miette::Result<String, AocError> {
    let stones = hailstone::parse(input)?;
    Ok(hailstone::crossings_in_area(&stones, low, high)?.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day24.txt");
        assert_eq!("2", process_with_area(input, 7, 27)?);
        Ok(())
    }

    #[rstest]
    #[case::crossing_on_the_corner_counts("0, 0, 0 @ 1, 1, 0\n4, 0, 0 @ 0, 1, 0", 0, 4, "1")]
    #[case::crossing_outside_the_area("0, 0, 0 @ 1, 1, 0\n4, 0, 0 @ 0, 1, 0", 0, 3, "0")]
    #[case::parallel_paths_never_cross("0, 0, 0 @ 1, 1, 0\n1, 0, 0 @ 2, 2, 0", -10, 10, "0")]
    #[case::crossed_in_the_past("0, 0, 0 @ -1, -1, 0\n4, 0, 0 @ 0, 1, 0", -10, 10, "0")]
    #[case::the_z_axis_is_ignored("0, 0, 0 @ 1, 1, 9\n4, 0, 0 @ 0, 1, -9", 0, 4, "1")]
    #[case::carriage_returns_and_spacing(
        "19, 13, 30 @ -2,  1, -2\r\n18, 19, 22 @ -1, -1, -2\r\n",
        7,
        27,
        "1"
    )]
    fn test_areas(
        #[case] input: &str,
        #[case] low: i128,
        #[case] high: i128,
        #[case] expected: &str,
    ) {
        assert_eq!(expected, process_with_area(input, low, high).unwrap());
    }

    #[test]
    fn the_real_test_area_stays_exact() {
        let input = concat!(
            "200000000000000, 200000000000000, 0 @ 1000, 1000, 0\n",
            "400000000000000, 200000000000000, 0 @ -1000, 1000, 0"
        );
        assert_eq!("1", process(input).unwrap());
        assert_eq!(
            "0",
            process_with_area(input, 200_000_000_000_000, 299_999_999_999_999).unwrap()
        );
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            "19, 13, 30",
            "19, 13, 30 @ -2, 1",
            "19, 13 @ -2, 1, -2",
            "x, 13, 30 @ -2, 1, -2",
            "19, 13, 30, 7 @ -2, 1, -2, 3",
        ] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
