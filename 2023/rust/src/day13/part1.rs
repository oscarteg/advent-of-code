use crate::day13::{custom_error::AocError, pattern::parse_patterns};

/// Sums the summary values of every pattern's exact mirror line.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let patterns = parse_patterns(input)?;
    let total: usize = patterns
        .iter()
        .map(|pattern| pattern.summarize(0))
        .sum::<Result<usize, AocError>>()?;
    Ok(total.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day13.txt");
        assert_eq!("405", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case(
        "#.##..##.\n..#.##.#.\n##......#\n##......#\n..#.##.#.\n..##..##.\n#.#.##.#.",
        "5"
    )]
    #[case(
        "#...##..#\n#....#..#\n..##..###\n#####.##.\n#####.##.\n..##..###\n#....#..#",
        "400"
    )]
    fn single_pattern_examples(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[test]
    fn vertical_mirror_at_the_edge() {
        let input = "##";
        assert_eq!("1", process(input).unwrap());
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "#.#\n.#x", "##\n#"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
