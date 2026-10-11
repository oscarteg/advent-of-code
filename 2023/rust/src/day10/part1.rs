use crate::day10::{custom_error::AocError, grid::Grid};

/// Finds how many steps along the main loop it takes to reach the tile
/// farthest from the starting position.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let grid = Grid::parse(input)?;
    let loop_tiles = grid.find_loop()?;
    Ok((loop_tiles.len() / 2).to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use rstest::rstest;

    #[rstest]
    #[case(include_str!("../../examples/day10.txt"), "4")]
    #[case(include_str!("../../examples/day10-2.txt"), "8")]
    fn test_process(#[case] input: &str, #[case] expected: &str) {
        assert_eq!(process(input).unwrap(), expected);
    }

    #[test]
    fn smallest_possible_loop() {
        let input = "S7\nLJ";
        assert_eq!(process(input).unwrap(), "2");
    }

    #[test]
    fn malformed_input_returns_errors() {
        for input in [
            "",
            ".....\n.X-7.\n.|.|.\n.L-J.\n.....",
            ".....\n.F-7.\n.|.|\n.L-J.\n.....",
            ".....\n..S..\n.....",
        ] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
