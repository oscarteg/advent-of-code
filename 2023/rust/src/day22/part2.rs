use crate::day22::{bricks::Stack, custom_error::AocError};

/// Sums how many bricks fall when each brick is removed in turn.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let stack = Stack::settle(input)?;
    Ok(stack.total_chain_reaction().to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case::official_example(include_str!("../../examples/day22.txt"), "7")]
    #[case::single_cube("0,0,1~0,0,1", "0")]
    #[case::column("0,0,1~0,0,1\n0,0,2~0,0,2\n0,0,3~0,0,3", "3")]
    #[case::shared_load("0,0,1~0,0,1\n1,0,1~1,0,1\n0,0,2~1,0,2", "0")]
    #[case::rests_on_the_highest_brick("0,0,1~0,0,1\n1,0,1~1,0,2\n0,0,3~1,0,3", "1")]
    #[case::cascade_stops_at_a_second_supporter(
        "0,0,1~0,0,1\n1,0,1~1,0,1\n0,0,2~1,0,2\n0,0,3~0,0,3",
        "1"
    )]
    fn chain_reactions(#[case] input: &str, #[case] expected: &str) -> miette::Result<()> {
        assert_eq!(expected, process(input)?);
        Ok(())
    }

    #[rstest]
    #[case::empty("")]
    #[case::missing_separator("1,0,1")]
    #[case::short_corner("1,0,1~1,0")]
    #[case::on_the_ground("1,0,0~1,2,0")]
    #[case::not_a_straight_line("0,0,1~1,1,1")]
    #[case::not_a_number("x,0,1~1,0,1")]
    fn malformed_input_is_rejected(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
