use crate::day08::{custom_error::AocError, network::Network};

const START: &str = "AAA";
const DESTINATION: &str = "ZZZ";

/// Counts the steps from `AAA` to `ZZZ`.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let network = Network::parse(input)?;
    let steps = network.steps(START, |name: &str| name == DESTINATION)?;
    Ok(steps.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(include_str!("../../examples/day08.txt"), "2")]
    #[case("LLR\n\nAAA = (BBB, BBB)\nBBB = (AAA, ZZZ)\nZZZ = (ZZZ, ZZZ)\n", "6")]
    #[case("LR\n\nAAA = (ZZZ, AAA)\nZZZ = (ZZZ, ZZZ)\n", "1")]
    fn walks_until_zzz(#[case] input: &str, #[case] expected: &str) -> miette::Result<()> {
        assert_eq!(expected, process(input)?);
        Ok(())
    }

    #[rstest]
    #[case("")]
    #[case("LR")]
    #[case("LX\n\nAAA = (AAA, ZZZ)\nZZZ = (ZZZ, ZZZ)\n")]
    #[case("LR\n\nAAA (BBB, CCC)\n")]
    #[case("LR\n\nAAA = BBB, CCC\n")]
    #[case("LR\n\nAAA = (BBB, CCC)\n")]
    #[case("LR\n\nBBB = (BBB, BBB)\n")]
    #[case("L\n\nAAA = (AAA, AAA)\nZZZ = (ZZZ, ZZZ)\n")]
    fn malformed_input_returns_errors(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
