use crate::day24::{custom_error::AocError, hailstone};

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let stones = hailstone::parse(input)?;
    Ok(hailstone::rock_position_sum(&stones)?.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    const THREE_HAILSTONES: &str = concat!(
        "12, 22, 31 @ -1, 0, 2\n",
        "8, 30, 34 @ 2, -3, 1\n",
        "13, 11, 45 @ 0, 5, -2\n"
    );

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day24.txt");
        assert_eq!("47", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case::three_hailstones_determine_the_rock(THREE_HAILSTONES, "60")]
    #[case::extra_hailstones_agree(
        concat!(
            "12, 22, 31 @ -1, 0, 2\n",
            "8, 30, 34 @ 2, -3, 1\n",
            "13, 11, 45 @ 0, 5, -2\n",
            "-18, 27, 58 @ 5, 1, -1\n"
        ),
        "60"
    )]
    #[case::the_order_does_not_matter(
        concat!(
            "20, 19, 15 @  1, -5, -3\n",
            "12, 31, 28 @ -1, -2, -1\n",
            "20, 25, 34 @ -2, -2, -4\n",
            "18, 19, 22 @ -1, -1, -2\n",
            "19, 13, 30 @ -2,  1, -2\n"
        ),
        "47"
    )]
    fn rock_positions(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[test]
    fn malformed_and_underdetermined_input_returns_errors() {
        for input in [
            "",
            "19, 13, 30 @ -2, 1, -2",
            "19, 13, 30 @ -2, 1, -2\n18, 19, 22 @ -1, -1, -2",
            "x, 13, 30 @ -2, 1, -2",
            "19, 13, 30 @ -2, 1",
        ] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
