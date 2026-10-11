use anyhow::{anyhow, bail, Context, Result};

fn parse_positions(input: &str) -> Result<Vec<u64>> {
    let positions = input
        .trim()
        .split(',')
        .map(|position| {
            position
                .trim()
                .parse::<u64>()
                .with_context(|| format!("invalid crab position: {position:?}"))
        })
        .collect::<Result<Vec<_>>>()?;
    if positions.is_empty() {
        bail!("at least one crab position is required");
    }
    Ok(positions)
}

fn linear_fuel(positions: &[u64], target: u64) -> Result<u64> {
    positions.iter().try_fold(0_u64, |fuel, position| {
        fuel.checked_add(position.abs_diff(target))
            .ok_or_else(|| anyhow!("linear fuel total overflows"))
    })
}

fn increasing_fuel(positions: &[u64], target: u64) -> Result<u64> {
    positions.iter().try_fold(0_u64, |fuel, position| {
        let distance = position.abs_diff(target);
        let (left, right) = if distance % 2 == 0 {
            (distance / 2, distance.checked_add(1))
        } else {
            (distance, distance.checked_add(1).map(|next| next / 2))
        };
        let steps = left
            .checked_mul(right.ok_or_else(|| anyhow!("fuel cost overflows"))?)
            .ok_or_else(|| anyhow!("fuel cost overflows"))?;
        fuel.checked_add(steps)
            .ok_or_else(|| anyhow!("increasing fuel total overflows"))
    })
}

/// Returns the minimum linear and increasing fuel costs.
pub fn solve(input: &str) -> Result<(u64, u64)> {
    let mut positions = parse_positions(input)?;
    positions.sort_unstable();

    let median = positions[positions.len() / 2];
    let part1 = linear_fuel(&positions, median)?;

    let (mut low, mut high) = (positions[0], positions[positions.len() - 1]);
    while low < high {
        let middle = low + (high - low) / 2;
        let middle_fuel = increasing_fuel(&positions, middle)?;
        let next_fuel = increasing_fuel(&positions, middle + 1)?;
        if middle_fuel <= next_fuel {
            high = middle;
        } else {
            low = middle + 1;
        }
    }

    Ok((part1, increasing_fuel(&positions, low)?))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn official_example() {
        let input = include_str!("../../input/day7-example.txt");
        assert_eq!(solve(input).unwrap(), (37, 168));
    }

    #[test]
    fn one_position_requires_no_fuel() {
        assert_eq!(solve("42").unwrap(), (0, 0));
    }

    #[test]
    fn identical_positions_require_no_fuel() {
        assert_eq!(solve("9,9,9,9").unwrap(), (0, 0));
    }

    #[test]
    fn supports_positions_at_numeric_boundaries() {
        assert_eq!(solve("0,1").unwrap(), (1, 1));
        assert!(solve("0,18446744073709551615").is_err());
    }

    #[test]
    fn rejects_empty_and_malformed_input() {
        for input in ["", " , ", "1,,2", "1,x,2", "-1,2"] {
            assert!(solve(input).is_err(), "accepted {input:?}");
        }
    }
}
