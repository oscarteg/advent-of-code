use crate::day08::{custom_error::AocError, network::Network};

/// Counts the steps until every ghost stands on a node ending in `Z`.
///
/// A ghost revisits its destination every walk length, so the first step the
/// ghosts share is the least common multiple of their walk lengths.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let network = Network::parse(input)?;
    let walks = network
        .starts()
        .map(|start| network.steps(start, |name: &str| name.ends_with('Z')))
        .collect::<Result<Vec<_>, _>>()?;
    let [first, rest @ ..] = walks.as_slice() else {
        return Err(AocError::InvalidInput("no start node ends in A"));
    };
    let steps = rest.iter().try_fold(*first, |steps, &walk| {
        least_common_multiple(steps, walk).ok_or(AocError::InvalidInput("step count overflows"))
    })?;
    Ok(steps.to_string())
}

fn least_common_multiple(a: u64, b: u64) -> Option<u64> {
    (a / greatest_common_divisor(a, b)).checked_mul(b)
}

fn greatest_common_divisor(a: u64, b: u64) -> u64 {
    match b {
        0 => a,
        _ => greatest_common_divisor(b, a % b),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(include_str!("../../examples/day08-2.txt"), "6")]
    #[case(include_str!("../../examples/day08.txt"), "2")]
    fn walks_until_the_ghosts_meet(
        #[case] input: &str,
        #[case] expected: &str,
    ) -> miette::Result<()> {
        assert_eq!(expected, process(input)?);
        Ok(())
    }

    #[rstest]
    #[case(2, 3, Some(6))]
    #[case(4, 6, Some(12))]
    #[case(7, 7, Some(7))]
    #[case(u64::MAX, 2, None)]
    fn least_common_multiple_reports_overflow(
        #[case] a: u64,
        #[case] b: u64,
        #[case] expected: Option<u64>,
    ) {
        assert_eq!(expected, least_common_multiple(a, b));
    }

    #[rstest]
    #[case("")]
    #[case("LX\n\n11A = (11Z, 11Z)\n11Z = (11Z, 11Z)\n")]
    #[case("LR\n\nBBB = (BBB, BBB)\n")]
    #[case("L\n\n11A = (11A, 11A)\n11Z = (11Z, 11Z)\n")]
    fn malformed_input_returns_errors(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
