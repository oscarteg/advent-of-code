use crate::day05::{almanac::Almanac, custom_error::AocError};

/// Finds the lowest location without enumerating the seed ranges.
#[tracing::instrument]
pub fn process(input: &str) -> miette::Result<String, AocError> {
    let almanac = Almanac::parse(input)?;
    let (pairs, remainder) = almanac.seeds.as_chunks::<2>();
    if !remainder.is_empty() {
        return Err(AocError::InvalidInput(
            "seed ranges require start/length pairs",
        ));
    }
    let seeds = pairs
        .iter()
        .map(|pair| {
            pair[0]
                .checked_add(pair[1])
                .map(|end| pair[0]..end)
                .ok_or(AocError::InvalidInput("seed range overflows"))
        })
        .collect::<Result<Vec<_>, _>>()?;
    let location = almanac
        .locations(seeds)
        .iter()
        .filter(|range| !range.is_empty())
        .map(|range| range.start)
        .min()
        .ok_or(AocError::InvalidInput("empty seed ranges"))?;
    Ok(location.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_process() -> miette::Result<()> {
        let input = include_str!("../../examples/day05.txt");
        assert_eq!("46", process(input)?);
        Ok(())
    }

    #[test]
    fn minimum_can_be_inside_a_seed_range() {
        let input = "seeds: 10 10\n\nseed-to-soil map:\n0 14 2";
        assert_eq!(process(input).unwrap(), "0");
    }

    #[test]
    fn enormous_ranges_are_not_enumerated() {
        let input = "seeds: 100 1000000000000\n\nseed-to-soil map:\n0 999999999999 1";
        assert_eq!(process(input).unwrap(), "0");
    }

    #[test]
    fn malformed_ranges_return_errors() {
        for input in ["seeds: 1", "seeds: 0 0", "seeds: 18446744073709551615 2"] {
            assert!(process(input).is_err(), "{input}");
        }
    }
}
