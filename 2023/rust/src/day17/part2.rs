use crate::day17::custom_error::AocError;
use crate::day17::grid::Grid;
use crate::day17::part1::least_heat_loss;

#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let grid = Grid::parse(input)?;
    let heat_loss =
        least_heat_loss(&grid, 4, 10).ok_or(AocError::InvalidInput("no path to the factory"))?;
    Ok(heat_loss.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day17.txt");
        assert_eq!("94", process(input)?);
        Ok(())
    }

    #[test]
    fn straight_corridor_forces_long_runs() -> miette::Result<()> {
        let input = include_str!("../../examples/day17-2.txt");
        assert_eq!("71", process(input)?);
        Ok(())
    }

    #[test]
    fn grid_too_small_for_minimum_run_has_no_path() {
        let input = "12\n34";
        assert!(process(input).is_err());
    }

    #[test]
    fn malformed_input_is_rejected() {
        for input in ["", "12\n1", "1a2\n345"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
