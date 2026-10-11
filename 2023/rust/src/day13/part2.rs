use crate::day13::{custom_error::AocError, pattern::parse_patterns};

/// Sums the summary values of every pattern's mirror line after fixing
/// exactly one smudge.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let patterns = parse_patterns(input)?;
    let total: usize = patterns
        .iter()
        .map(|pattern| pattern.summarize(1))
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
        assert_eq!("400", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case(
        "#.##..##.\n..#.##.#.\n##......#\n##......#\n..#.##.#.\n..##..##.\n#.#.##.#.",
        "300"
    )]
    #[case(
        "#...##..#\n#....#..#\n..##..###\n#####.##.\n#####.##.\n..##..###\n#....#..#",
        "100"
    )]
    fn single_pattern_examples(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[test]
    fn pattern_with_no_single_smudge_fix_returns_error() {
        let input = "#.\n.#";
        assert!(process(input).is_err());
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", "#.#\n.#x", "##\n#"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
