use crate::day05::{almanac::Almanac, custom_error::AocError};

/// Finds the lowest location of the individually listed seeds.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let almanac = Almanac::parse(input)?;
    let location = almanac
        .seeds
        .iter()
        .map(|&seed| almanac.location(seed))
        .min()
        .ok_or(AocError::InvalidInput("empty seed list"))?;
    Ok(location.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day05.txt");
        assert_eq!("35", process(input)?);
        Ok(())
    }

    #[test]
    fn source_end_is_exclusive() {
        let input = "seeds: 12\n\nseed-to-soil map:\n100 10 2";
        assert_eq!(process(input).unwrap(), "12");
    }

    #[test]
    fn mapped_values_are_not_mapped_twice_in_one_stage() {
        let input = "seeds: 10\n\nseed-to-soil map:\n20 10 1\n0 20 1";
        assert_eq!(process(input).unwrap(), "20");
    }
}
