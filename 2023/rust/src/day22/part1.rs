use crate::day22::{bricks::Stack, custom_error::AocError};

/// Counts the bricks that can be removed without anything falling.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let stack = Stack::settle(input)?;
    Ok(stack.safe_to_disintegrate().to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case::official_example(include_str!("../../examples/day22.txt"), "5")]
    #[case::single_cube("0,0,1~0,0,1", "1")]
    #[case::cube_falls_to_the_ground("0,0,9~0,0,9", "1")]
    #[case::column("0,0,1~0,0,1\n0,0,2~0,0,2\n0,0,3~0,0,3", "1")]
    #[case::shared_load("0,0,1~0,0,1\n1,0,1~1,0,1\n0,0,2~1,0,2", "3")]
    #[case::rests_on_the_highest_brick("0,0,1~0,0,1\n1,0,1~1,0,2\n0,0,3~1,0,3", "2")]
    #[case::snapshot_order_is_irrelevant(
        "1,1,8~1,1,9\n0,1,6~2,1,6\n2,0,5~2,2,5\n0,0,4~0,2,4\n0,2,3~2,2,3\n0,0,2~2,0,2\n1,0,1~1,2,1",
        "5"
    )]
    fn settled_stacks(#[case] input: &str, #[case] expected: &str) -> miette::Result<()> {
        assert_eq!(expected, process(input)?);
        Ok(())
    }

    #[rstest]
    #[case::empty("")]
    #[case::missing_separator("1,0,1")]
    #[case::short_corner("1,0,1~1,0")]
    #[case::long_corner("1,0,1~1,0,1,1")]
    #[case::on_the_ground("1,0,0~1,2,0")]
    #[case::not_a_straight_line("0,0,1~1,1,1")]
    #[case::not_a_number("x,0,1~1,0,1")]
    #[case::negative_coordinate("0,0,1~0,0,-1")]
    fn malformed_input_is_rejected(#[case] input: &str) {
        assert!(process(input).is_err(), "{input}");
    }
}
