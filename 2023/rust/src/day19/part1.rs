use crate::day19::{custom_error::AocError, workflows::System};

/// Sums the ratings of every part the workflows accept.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let system = System::parse(input)?;
    let mut total = 0;
    for &part in &system.parts {
        if system.accepts(part)? {
            total += part.total_rating();
        }
    }
    Ok(total.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day19.txt");
        assert_eq!("19114", process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("in{A}\n\n{x=1,m=2,a=3,s=4}", "10")]
    #[case("in{R}\n\n{x=1,m=2,a=3,s=4}", "0")]
    #[case("in{s<2:hop,R}\nhop{A}\n\n{x=1,m=1,a=1,s=1}", "4")]
    #[case("in{x<2:A,R}\n\n{x=1,m=1,a=1,s=1}", "4")]
    #[case("in{x<2:A,R}\n\n{x=2,m=1,a=1,s=1}", "0")]
    #[case("in{x>5:R,x>1:A,R}\n\n{x=10,m=1,a=1,s=1}\n{x=3,m=1,a=1,s=1}", "6")]
    fn small_systems(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(expected, process(input).unwrap());
    }

    #[rstest]
    #[case("")]
    #[case("in{A}")]
    #[case("{x=1,m=1,a=1,s=1}")]
    #[case("px{A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in A\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{A}\nin{R}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{x<1:A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{A,x<1:R}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{q<1:A,A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{x=1:A,A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{x<1:gone,A}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{loop}\nloop{in}\n\n{x=1,m=1,a=1,s=1}")]
    #[case("in{A}\n\n{x=1,m=1,a=1}")]
    #[case("in{A}\n\n{x=1,m=1,a=1,s=-1}")]
    #[case("in{A}\n\nx=1,m=1,a=1,s=1")]
    fn malformed_input_is_rejected(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
