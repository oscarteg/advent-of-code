use crate::day19::{custom_error::AocError, workflows::System};

/// Counts the distinct rating combinations the workflows accept.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let combinations = System::parse(input)?.accepted_combinations()?;
    Ok(combinations.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    /// Every rating combination of `1..=4000` in four categories.
    const ALL: u64 = 4000 * 4000 * 4000 * 4000;
    /// One rating in one category, the other three unconstrained.
    const ONE: u64 = 4000 * 4000 * 4000;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day19.txt");
        assert_eq!("167409079868000", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("in{A}\n\n{x=1,m=1,a=1,s=1}", ALL)]
    #[case("in{R}\n\n{x=1,m=1,a=1,s=1}", 0)]
    #[case("in{x<2001:A,A}\n\n{x=1,m=1,a=1,s=1}", ALL)]
    #[case("in{x<2:A,R}\n\n{x=1,m=1,a=1,s=1}", ONE)]
    #[case("in{x>3999:A,R}\n\n{x=1,m=1,a=1,s=1}", ONE)]
    #[case("in{x<1:A,R}\n\n{x=1,m=1,a=1,s=1}", 0)]
    #[case("in{x>4000:A,R}\n\n{x=1,m=1,a=1,s=1}", 0)]
    #[case("in{x>10:A,x>5:A,R}\n\n{x=1,m=1,a=1,s=1}", 3995 * ONE)]
    #[case("in{s<2:hop,R}\nhop{A}\n\n{x=1,m=1,a=1,s=1}", ONE)]
    fn small_systems(#[case] input: &str, #[case] expected: u64) {
        assert_eq!(expected.to_string(), process(input).unwrap());
    }

    #[rstest]
    #[case("")]
    #[case("in{A}")]
    #[case("in{x<1:A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{x<1:gone,A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{loop}\nloop{in}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{x<}\n\n{x=1,m=1,a=1,s=1}")]
    fn malformed_input_is_rejected(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
