use crate::day10::{
    custom_error::AocError,
    grid::{enclosed_area, Grid},
};

/// Counts the tiles enclosed by the main loop.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let grid = Grid::parse(input)?;
    let loop_tiles = grid.find_loop()?;
    Ok(enclosed_area(&loop_tiles).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(include_str!("../../examples/day10-3.txt"), "4")]
    #[case(include_str!("../../examples/day10-4.txt"), "8")]
    fn test_process(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(process(input).unwrap(), expected);
    }

    #[test]
    fn smallest_possible_loop_has_no_enclosed_tile() {
        let input = "S7\nLJ";
        assert_eq!(process(input).unwrap(), "0");
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in ["", ".....\n.X-7.\n.|.|.\n.L-J.\n....."] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
